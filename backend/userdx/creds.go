package userdx

import (
	"context"
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

var (
	ErrCodePasswordReuse = "idx.err.passwordReuse"
	ErrPasswordReuse     = errors.New("password reused")

	ErrCodePasswordExpired = "idx.err.passwordExpired"
	ErrPasswordExpired     = errors.New("password expired")

	ErrCodeTooManyFailedAttempts = "idx.err.tooManyFailedAttempts"
	ErrTooManyFailedAttempts     = "too many failed login attempts"

	ErrCodeInvalidCreds = "idx.err.invalidCreds"
	ErrInvalidCreds     = "invalid credential provided"
)

type SecretStorage struct {
	hasher     core.Hasher
	pwPolicy   map[core.AuthEntity]*core.CredentialPolicy
	policyLock sync.RWMutex
}

func NewCredentialStorage(
	hasher core.Hasher) core.SecretStorage {
	return &SecretStorage{
		hasher:   hasher,
		pwPolicy: make(map[core.AuthEntity]*core.CredentialPolicy),
	}
}

func (pcs *SecretStorage) CreatePassword(
	gtx context.Context, creds *core.Creds) error {

	policy, err := pcs.CredentialPolicy(gtx, creds.Type)
	if err != nil {
		return errx.Wrap(err)
	}

	if err = policy.MatchPattern(creds.Password); err != nil {
		return errx.Wrap(err)
	}

	hash, err := pcs.hasher.Hash(creds.Password)
	if err != nil {
		return errx.Wrap(err)
	}
	const query = `
		INSERT INTO credential (
			unique_name,
			item_type,
			password_hash,
			created_at,
			retries,
		) VALUES (
			$1,
			$2,
			$3
		) 
	`

	_, err = pg.Conn().ExecContext(
		gtx, query, creds.UniqueName, creds.Type, hash)
	if err != nil {
		return errx.Errf(err,
			"failed to update password hash for '%s - %s' in DB",
			creds.UniqueName, creds.Type)
	}
	return nil
}

func (pcs *SecretStorage) UpdatePassword(
	gtx context.Context, creds *core.Creds) error {

	// check if password matches the policy
	polocy, err := pcs.CredentialPolicy(gtx, creds.Type)
	if err != nil {
		return errx.Wrap(err)
	}

	if err := polocy.MatchPattern(creds.Password); err != nil {
		return errx.Errf(err, "passwrod does not meet complexity requirement")
	}

	hash, err := pcs.hasher.Hash(creds.Password)
	if err != nil {
		return errx.Wrap(err)
	}

	cur, err := pcs.getStoredCreds(gtx, creds)
	if err != nil {
		return errx.Wrap(err)
	}

	if slices.Contains(cur.PrevPasswords, hash) {
		return errx.Errfx(
			ErrPasswordReuse, ErrCodePasswordReuse, "cannot reuse passwords")
	}

	cur.PrevPasswords = append(cur.PrevPasswords, hash)
	pp := cur.PrevPasswords
	if len(cur.PrevPasswords) > polocy.MaxReuse {
		pp = cur.PrevPasswords[1:]
	}

	// TODO - Check if the new password is already is in prevPasswords...

	// TODO - how to handle array
	const query = `
			UPDATE credential SET
				password_hash = $1,
				created_at = $2,
				retries = 0,
				prev_passwords = $3
			WHERE 
				unique_name = $4,
				item_id = $5
			;
		`
	_, err = pg.Conn().ExecContext(
		gtx,
		query,
		hash,
		time.Now(),
		pp,
		creds.UniqueName,
		creds.Type,
	)
	if err != nil {
		return errx.Errf(err,
			"failed to update password hash for '%s (%s)' in DB",
			creds.UniqueName, creds.Type)
	}
	return nil
}

func (pcs *SecretStorage) Authenticate(
	gtx context.Context, in *core.Creds) error {

	secret, err := pcs.getStoredCreds(gtx, in)
	if err != nil {
		return errx.Wrap(err)
	}

	policy, err := pcs.CredentialPolicy(gtx, secret.Type)
	if err != nil {
		return errx.Wrap(err)
	}

	// Check if password has expired
	if secret.CreatedOn.Add(policy.Expiry).After(time.Now()) {
		return errx.Errfx(
			ErrPasswordExpired,
			ErrCodePasswordExpired,
			"password has expired, please reset password")
	}

	if secret.NumFailedAuth > policy.MaxRetries {

		resetInterval := time.Hour * time.Duration(policy.RetryResetDays*24)
		resetTime := secret.LastFailedOn.Add(resetInterval)
		if resetTime.After(time.Now()) {
			return errx.Errfx(
				ErrPasswordExpired,
				ErrCodePasswordExpired,
				"password has expired, please reset password")
		}
		//TODO  Reset the num_failed_auth
		query := `
			UPDATE credential SET 
				num_failed_auth = num_failed_auth + 1,
				last_failed_on = NOW()
		`
		if _, err = pg.Conn().ExecContext(gtx, query); err != nil {
			return errx.Errf(err, "invalid credentials: '%s (%s)', "+
				"failed to update failure count", in.UniqueName, in.Type)
		}
	}

	if err = pcs.hasher.Verify(in.Password, secret.PasswordHash); err != nil {

		query := `
			UPDATE credential SET 
				num_failed_auth = num_failed_auth + 1,
				last_failed_on = NOW()
		`
		if _, err = pg.Conn().ExecContext(gtx, query); err != nil {
			return errx.Errf(err, "invalid credentials: '%s (%s)', "+
				"failed to update failure count", in.UniqueName, in.Type)
		}

		// TODO - execute the query

		return errx.Errfx(err, ErrCodeInvalidCreds,
			"invalid credentials: '%s (%s)'",
			in.UniqueName, in.Type)
	}

	return nil
}

func (pcs *SecretStorage) getStoredCreds(
	gtx context.Context, givenCreds *core.Creds) (*core.Secret, error) {
	const query = `
	SELECT * 
	FROM credential
	WHERE 
		unique_name = $1 AND
		item_type = $2
	`
	var creds core.Secret
	err := pg.Conn().GetContext(
		gtx, &creds, query, givenCreds.UniqueName, givenCreds.Type)
	if err != nil {
		return nil, errx.Errf(err,
			"failed to get credential info from DB for '%s'",
			creds.UniqueName)
	}
	return &creds, nil
}

func (pcs *SecretStorage) StoreToken(
	gtx context.Context, token *core.Token) error {
	const query = `
		INSERT INTO idx_token(
			token,
			unique_name,
			assoc_type,
			operation		
		) VALUES (
			:token,
			:unique_name,
			:assoc_type,
			:operation
		) 
	`

	if _, err := pg.Conn().NamedExecContext(gtx, query, token); err != nil {
		return errx.Errf(err, "failed to create token for '%s:%s:%s",
			token.AssocType, token.Operation, token.UniqueName)
	}
	return nil
}

func (pcs *SecretStorage) VerifyToken(
	gtx context.Context, un, operation, token string) error {
	const query = `
		SELECT EXISTS( 
			SELECT 1 
			FROM idx_token 
			WHERE 
				token = $1 AND 
				unique_name = $2 AND 
				operation = $3
			)
	`
	exists := false
	err := pg.Conn().GetContext(gtx, &exists, query, token, un, operation)
	if err != nil {
		return errx.Errf(err,
			"failed to verify toke for %s (%s)", un, operation)
	}

	const dquery = `
		DELETE FROM idx_token 
		WHERE 
			token = $1 AND 
			unique_name = $2 AND 
			operation = $3`
	_, err = pg.Conn().ExecContext(gtx, dquery, token, un, operation)
	if err != nil {
		id := un + ":" + operation
		log.Error().Str("token", id).Msg("failed to delete verified token")
	}

	return nil
}

func (pcs *SecretStorage) CredentialPolicy(
	gtx context.Context,
	credType core.AuthEntity) (*core.CredentialPolicy, error) {
	pcs.policyLock.RLock()
	defer pcs.policyLock.RUnlock()

	if policy, found := pcs.pwPolicy[credType]; found {
		return policy, nil
	}

	const query = `SELECT * FROM credential_policy WHERE item_type = $1`
	var policy core.CredentialPolicy
	err := pg.Conn().SelectContext(gtx, &policy, query, credType)
	if err != nil {
		return nil, errx.Errf(err, "failed to retrieve cred policy from DB")
	}

	pcs.pwPolicy[credType] = &policy
	return &policy, nil
}

func (pcs *SecretStorage) SetCredentialPolicy(
	gtx context.Context,
	cp *core.CredentialPolicy) error {

	const query = `INSERT INTO 
		credential_policy (
			item_type,
			pattern,
			expiry,
			max_retries,
			max_reuse	
		) VALUES (
			:item_type,
			:pattern,
			:expiry,
			:max_retries,
			:max_reuse
		) ON CONFLICT(item_type) DO UPDATE SET 
		 	item_type = EXCLUDED.item_type,
			pattern = EXCLUDED.pattern,
			expiry = EXCLUDED.expiry,
			max_retries = EXCLUDED.max_retries,
			max_reuse = EXCLUDED.max_reuse
		;`

	if _, err := pg.Conn().NamedExecContext(gtx, query, &cp); err != nil {
		return errx.Errf(err, "failed to create/update creds policy")
	}

	pcs.policyLock.Lock()
	defer pcs.policyLock.Unlock()
	pcs.pwPolicy[cp.ItemType] = cp
	return nil
}
