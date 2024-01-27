package pg

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

type UserStorage struct {
	gd *pg.GetterDeleter
}

func NewUserStorage(gd *pg.GetterDeleter) core.UserStorage {
	return &UserStorage{
		gd: gd,
	}
}

func (pgu *UserStorage) Save(
	gtx context.Context, user *core.User) (err error) {

	query := `
		INSERT INTO idx_user (
			created_by,
			updated_by,
			user_id,
			email,
			auth,
			first_name,
			last_name,
			title,
			props			
		) VALUES (
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
				created_by = EXCLUDED.created_by,
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

func (pgu *UserStorage) Update(
	gtx context.Context, user *core.User) error {

	query := `
		UPDATE idx_user SET
			created_by = :created_by,
			updated_by = :updated_by,
			user_id = :user_id,
			email = :email,
			auth = :auth,
			first_name = :first_name,
			last_name = :last_name,
			title = :title,
			props = :props
		WHERE id = :id	
	`
	if _, err := pg.Conn().NamedExecContext(gtx, query, user); err != nil {
		return errx.Errf(
			err, "failed to insert user '%s' to database", user.Id())
	}

	return nil
}

func (pgu *UserStorage) GetOne(
	gtx context.Context, id string) (*core.User, error) {
	var user core.User
	err := pgu.gd.GetOne(gtx, "idx_user", "user_id", id, &user)
	if err != nil {
		return nil, err
	}

	// get groups and permissions

	return &user, nil
}

func (pgu *UserStorage) SetState(
	gtx context.Context, id string, state core.UserState) error {
	return nil
}

func (pgu *UserStorage) Remove(gtx context.Context, id string) error {
	return nil
}

func (pgu *UserStorage) Get(
	gtx context.Context, params data.CommonParams) ([]*core.User, error) {
	return nil, nil
}

func (pgu *UserStorage) AddToGroup(
	gtx context.Context, userId, groupId int) error {
	return nil
}

func (pgu *UserStorage) RemoveFromGroup(
	gtx context.Context, userId, groupId int) error {
	return nil
}
