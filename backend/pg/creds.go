package pg

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

type CredentialStorage struct {
	hasher core.Hasher
}

func NewCredentialStorage(hasher core.Hasher) core.CredentialStorage {
	return &CredentialStorage{
		hasher: hasher,
	}
}

func (pcs *CredentialStorage) SetPassword(
	gtx context.Context, itemType, id, password string) error {
	hash, err := pcs.hasher.Hash(password)
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
	_, err = pg.Conn().ExecContext(gtx, query, id, itemType, hash)
	if err != nil {
		return errx.Errf(err,
			"failed to update password hash for '%s - %s' in DB",
			id, itemType)
	}
	return nil
}

func (pcs *CredentialStorage) UpdatePassword(
	gtx context.Context, itemType, id, oldPw, newPw string) error {

	if err := pcs.Verify(gtx, itemType, id, oldPw); err != nil {
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
	_, err = pg.Conn().ExecContext(gtx, query, id, itemType, hash)
	if err != nil {
		return errx.Errf(err,
			"failed to update password hash for '%s (%s)' in DB",
			id, itemType)
	}
	return nil
}

func (pcs *CredentialStorage) Verify(
	gtx context.Context, itemType, id, password string) error {
	const query = `
		SELECT password_hash 
		FROM credential
		WHERE 
			id = $1,
			item_type = $2
	`
	hash := ""
	err := pg.Conn().GetContext(gtx, &hash, query, id, itemType)
	if err != nil {
		return errx.Errf(err, "failed to get password info from DB")
	}

	ok, err := pcs.hasher.Verify(password, hash)
	if err != nil {
		return err
	}
	if !ok {
		return errx.Errf(err,
			"failed to verify password for '%s (%s)'", id, itemType)
	}

	return nil
}
