package pg

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

type SecretStorage struct {
	hasher core.Hasher
}

func NewCredentialStorage(
	hasher core.Hasher) core.SecretStorage {
	return &SecretStorage{
		hasher: hasher,
	}
}

func (pcs *SecretStorage) SetPassword(
	gtx context.Context, creds *core.Creds) error {
	hash, err := pcs.hasher.Hash(creds.Password)
	if err != nil {
		return err
	}
	const query = `
		INSERT INTO credential (
			id,
			item_type,
			password_hash
		) VALUES (
			$1,
			$2,
			$3
		) ON CONFLICT(id) DO
			UPDATE SET password_hash = EXCLUDED.password_hash;
	
	`
	_, err = pg.Conn().ExecContext(gtx, query, creds.Id, creds.Type, hash)
	if err != nil {
		return errx.Errf(err,
			"failed to update password hash for '%s - %s' in DB",
			creds.Id, creds.Type)
	}
	return nil
}

func (pcs *SecretStorage) UpdatePassword(
	gtx context.Context, creds *core.Creds, newPw string) error {

	if err := pcs.Verify(gtx, creds); err != nil {
		return err
	}

	hash, err := pcs.hasher.Hash(newPw)
	if err != nil {
		return err
	}
	const query = `
			UPDATE credential SET
				password_hash = $1
			WHERE 
				id = $2,
				item_id = $3
			;
		`
	_, err = pg.Conn().ExecContext(gtx, query, hash, creds.Id, creds.Type)
	if err != nil {
		return errx.Errf(err,
			"failed to update password hash for '%s (%s)' in DB",
			creds.Id, creds.Type)
	}
	return nil
}

func (pcs *SecretStorage) Verify(gtx context.Context, creds *core.Creds) error {
	const query = `
		SELECT password_hash 
		FROM credential
		WHERE 
			id = $1,
			item_type = $2
	`
	hash := ""
	err := pg.Conn().GetContext(gtx, &hash, query, creds.Id, creds.Type)
	if err != nil {
		return errx.Errf(err, "failed to get password info from DB")
	}

	ok, err := pcs.hasher.Verify(creds.Password, hash)
	if err != nil {
		return err
	}
	if !ok {
		return errx.Errf(err,
			"failed to verify password for '%s (%s)'", creds.Id, creds.Type)
	}

	return nil
}

func (pcs *SecretStorage) StoreToken(
	gtx context.Context, token *core.Token) error {
	const query = `
		INSERT INTO idx_token(
			token,
			id,
			assoc_type,
			operation		
		) VALUES (
			:token,
			:id,
			:assoc_type,
			:operation					
		) 
	`

	if _, err := pg.Conn().NamedExecContext(gtx, query, token); err != nil {
		return errx.Errf(err, "failed to create token for '%s:%s:%s",
			token.AssocType, token.Operation, token.Id)
	}
	return nil
}

func (pcs *SecretStorage) VerifyToken(
	gtx context.Context, token, id, operation string) error {
	const query = `
		SELECT EXISTS( 
			SELECT 1 
			FROM idx_token 
			WHERE 
				token = $1 AND 
				id = $2 AND 
				operation = $3
			)
	`
	exists := false
	err := pg.Conn().GetContext(gtx, &exists, query, token, id, operation)
	if err != nil {
		return errx.Errf(err,
			"failed to verify toke for %s (%s)", id, operation)
	}

	const dquery = `
		DELETE FROM idx_token 
		WHERE 
			token = $1 AND 
			id = $2 AND 
			operation = $3`
	_, err = pg.Conn().ExecContext(gtx, dquery, token, id, operation)
	if err != nil {
		id := id + ":" + operation
		slog.Error("failed to delete verified token", "tokenId", id)
	}

	return nil
}

func NewToken(assocType, id, operation string) *core.Token {
	return &core.Token{
		Token:     uuid.NewString(),
		Id:        id,
		AssocType: assocType,
		Operation: operation,
	}
}
