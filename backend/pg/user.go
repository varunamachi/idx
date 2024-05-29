package pg

import (
	"context"
	"time"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

type userPgStorage struct {
	gd data.GetterDeleter
}

func NewUserStorage(gd data.GetterDeleter) core.UserStorage {
	return &userPgStorage{
		gd: gd,
	}
}

func (pgu *userPgStorage) Save(
	gtx context.Context, user *core.User) (int64, error) {

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
		) ON CONFLICT (id) DO UPDATE SET
				created_by = EXCLUDED.created_by,
				updated_by = EXCLUDED.updated_by,
				user_id = EXCLUDED.user_id,
				email = EXCLUDED.email,
				auth = EXCLUDED.auth,
				first_name = EXCLUDED.first_name,
				last_name = EXCLUDED.last_name,
				title = EXCLUDED.title,
				props = EXCLUDED.props
		RETURNING id;
	`

	stmt, err := pg.Conn().PrepareNamed(query)
	if err != nil {
		return -1, errx.Errf(err, "failed to prepare query to save user")
	}

	var id int64
	if err = stmt.GetContext(gtx, &id, user); err != nil {
		return -1, errx.Errf(
			err, "failed to insert user '%s' to database", user.Id())
	}
	return id, nil

	// if _, err := pg.Conn().NamedExecContext(gtx, query, user); err != nil {
	// 	return -1, errx.Errf(
	// 		err, "failed to insert user '%s' to database", user.Id())
	// }

	// return 0, nil
}

func (pgu *userPgStorage) Update(
	gtx context.Context, user *core.User) error {

	user.UpdatedAt = time.Now()
	query := `
		UPDATE idx_user SET
			updated_by = :updated_by,
			updated_at = :updated_at,
			user_id = :user_id,
			email = :email,
			auth = :auth,
			state = :state,
			first_name = :first_name,
			last_name = :last_name,
			title = :title,
			props = :props
		WHERE id = :id	
	`
	if _, err := pg.Conn().NamedExecContext(gtx, query, user); err != nil {
		return errx.Errf(
			err, "failed to update user '%s' to database", user.Id())
	}

	return nil
}

func (pgu *userPgStorage) GetOne(
	gtx context.Context, id int64) (*core.User, error) {
	var user core.User
	err := pgu.gd.GetOne(gtx, "idx_user", "id", id, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pgu *userPgStorage) GetByUserId(
	gtx context.Context, id string) (*core.User, error) {
	var user core.User
	err := pgu.gd.GetOne(gtx, "idx_user", "user_id", id, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pgu *userPgStorage) SetState(
	gtx context.Context, id int64, state core.UserState) error {
	query := `
		UPDATE idx_user SET
			state = $1,
		WHERE id = $2	
	`

	_, err := pg.Conn().ExecContext(gtx, query, id, state)
	if err != nil {
		return errx.Errf(
			err, "failed to insert user '%s' to database", id)
	}

	return nil
}

func (pgu *userPgStorage) Remove(gtx context.Context, id int64) error {
	query := `DELETE FROM idx_user WHERE id = $2`

	_, err := pg.Conn().ExecContext(gtx, query, id)
	if err != nil {
		return errx.Errf(
			err, "failed to insert user '%s' to database", id)
	}

	return nil
}

func (pgu *userPgStorage) Get(
	gtx context.Context, params *data.CommonParams) ([]*core.User, error) {

	out := make([]*core.User, 0, params.PageSize)

	if err := pgu.gd.Get(gtx, "idx_user", params, &out); err != nil {
		return nil, err
	}

	// Note: Get permissions per user per service on demand

	return out, nil
}

func (pgu *userPgStorage) Exists(
	gtx context.Context, userId string) (bool, error) {
	return pgu.gd.Exists(gtx, "idx_user", "user_id", userId)
}

func (pgu *userPgStorage) Count(
	gtx context.Context, filter *data.Filter) (int64, error) {
	return pgu.gd.Count(gtx, "idx_user", filter)
}
