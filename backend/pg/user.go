package pg

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

type UserStorage struct {
	gd data.GetterDeleter
}

func NewUserStorage(gd data.GetterDeleter) core.UserStorage {
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

func (pgu *UserStorage) GetOne(
	gtx context.Context, id int) (*core.User, error) {
	var user core.User
	err := pgu.gd.GetOne(gtx, "idx_user", "id", id, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pgu *UserStorage) SetState(
	gtx context.Context, id int, state core.UserState) error {
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

func (pgu *UserStorage) Remove(gtx context.Context, id int) error {
	query := `DELETE FROM idx_user WHERE id = $2`

	_, err := pg.Conn().ExecContext(gtx, query, id)
	if err != nil {
		return errx.Errf(
			err, "failed to insert user '%s' to database", id)
	}

	return nil
}

func (pgu *UserStorage) Get(
	gtx context.Context, params data.CommonParams) ([]*core.User, error) {

	out := make([]*core.User, 0, params.PageSize)

	if err := pgu.gd.Get(gtx, "idx_user", &params, &out); err != nil {
		return nil, err
	}

	// TODO:  What about permissions

	return out, nil
}

func (pgu *UserStorage) AddToGroup(
	gtx context.Context, userId, groupId int) error {

	query := `
		INSERT INTO user_to_group (
			user_id, 
			group_id
		)
		VALUES (
			$1, 
			$2 
		)	
	`
	_, err := pg.Conn().ExecContext(gtx, query, userId, groupId)
	if err != nil {
		return errx.Errf(
			err, "failed to add user '%s' to group '%s'", userId, groupId)
	}
	return nil
}

func (pgu *UserStorage) RemoveFromGroup(
	gtx context.Context, userId, groupId int) error {
	query := `
		DELETE FROM user_to_group WHERE user_id = $1 AND group_id = $2 	
	`
	_, err := pg.Conn().ExecContext(gtx, query, userId, groupId)
	if err != nil {
		return errx.Errf(
			err, "failed to remove user '%s' from group '%s'", userId, groupId)
	}
	return nil
}

func (pgu *UserStorage) GetPermissionForService(
	gtx context.Context, userId, serviceId int) ([]string, error) {
	query := `
		SELECT
			perm_id
		FROM group_to_perm g2p
		JOIN idx_group g ON g.id = g2p.group_id
		JOIN user_to_group u2g ON g.id = u2g.group_id
		JOIN idx_user u ON u.id = u2g.id
		JOIN service_to_group s2g ON s2g.group_id = g.id
		WHERE u.user_id = $1 AND s2g.service_id = $2
	`

	perms := make([]string, 0, 50)
	err := pg.Conn().SelectContext(gtx, &perms, query, userId, serviceId)
	if err != nil {
		return nil,
			errx.Errf(
				err,
				"failed to get permissions of user '%s' for service '%s'",
				userId, serviceId)

	}
	return perms, nil
}
