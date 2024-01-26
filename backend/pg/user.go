package pg

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

type PgUserStorage struct {
}

func (pgu *PgUserStorage) Save(
	gtx context.Context, user *core.User) (err error) {

	query := `
		INSERT INTO idx_user (
			id,
			created_at,
			created_by,
			updated_at,
			updated_by,
			user_id,
			email,
			auth,
			first_name,
			last_name,
			title,
			props			
		) VALUES (
			:id,
			:created_at,
			:created_by,
			:updated_at,
			:updated_by,
			:user_id,
			:email,
			:auth,
			:first_name,
			:last_name,
			:title,
			:props	
		) ON CONFLICT (user_id) DO UPDATE SET
				created_at = EXCLUDED.created_at,
				created_by = EXCLUDED.created_by,
				updated_at = EXCLUDED.updated_at,
				updated_by = EXCLUDED.updated_by,
				user_id = EXCLUDED.user_id,
				email = EXCLUDED.email,
				auth = EXCLUDED.auth,
				first_name = EXCLUDED.first_name,
				last_name = EXCLUDED.last_name,
				title = EXCLUDED.title,
				props = EXCLUDED.props
	`

	if _, err := pg.Conn().NamedExecContext(gtx, query, user); err != nil {
		return errx.Errf(
			err, "failed to insert user '%s' to database", user.Id())
	}

	return nil
}

func (pgu *PgUserStorage) Update(
	gtx context.Context, user *core.User) (err error) {
	return nil
}

func (pgu *PgUserStorage) GetOne(
	gtx context.Context, id string) (*core.User, error) {
	return nil, nil
}

func (pgu *PgUserStorage) SetState(
	gtx context.Context, id string, state core.UserState) error {
	return nil
}

func (pgu *PgUserStorage) Remove(gtx context.Context, id string) error {
	return nil
}

func (pgu *PgUserStorage) Get(
	gtx context.Context, params data.CommonParams) ([]*core.User, error) {
	return nil, nil
}

func (pgu *PgUserStorage) AddToGroup(
	gtx context.Context, userId, groupId int) error {
	return nil
}

func (pgu *PgUserStorage) RemoveFromGroup(
	gtx context.Context, userId, groupId int) error {
	return nil
}
