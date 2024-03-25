package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
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

func (pgu *UserStorage) GetOne(
	gtx context.Context, id int64) (*core.User, error) {
	var user core.User
	err := pgu.gd.GetOne(gtx, "idx_user", "id", id, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pgu *UserStorage) GetByUserId(
	gtx context.Context, id string) (*core.User, error) {
	var user core.User
	err := pgu.gd.GetOne(gtx, "idx_user", "user_id", id, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pgu *UserStorage) SetState(
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

func (pgu *UserStorage) Remove(gtx context.Context, id int64) error {
	query := `DELETE FROM idx_user WHERE id = $2`

	_, err := pg.Conn().ExecContext(gtx, query, id)
	if err != nil {
		return errx.Errf(
			err, "failed to insert user '%s' to database", id)
	}

	return nil
}

func (pgu *UserStorage) Get(
	gtx context.Context, params *data.CommonParams) ([]*core.User, error) {

	out := make([]*core.User, 0, params.PageSize)

	if err := pgu.gd.Get(gtx, "idx_user", params, &out); err != nil {
		return nil, err
	}

	// Note: Get permissions per user per service on demand

	return out, nil
}

func (pgu *UserStorage) AddToGroups(
	gtx context.Context, userId int64, groupIds ...int64) error {

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

	tx, err := pg.Conn().BeginTxx(gtx, &sql.TxOptions{})
	if err != nil {
		return errx.Errf(err, "failed to initilize DB transaction")
	}
	ef := func(err error, fmtStr string, args ...any) error {
		if e := tx.Rollback(); e != nil {
			log.Error().Err(err).
				Msg("transaction rollback failed for fake users table")
		}
		return errx.Errf(err, fmtStr, args...)
	}

	for _, gid := range groupIds {
		_, err := tx.ExecContext(gtx, query, userId, gid)
		if err != nil {
			return ef(
				err, "failed to add user '%s' to group '%s'", userId, gid)
		}
	}

	if err := tx.Commit(); err != nil {
		return ef(err, "failed commit DB transaction: user id '%s'", userId)
	}

	return nil
}

func (pgu *UserStorage) RemoveFromGroup(
	gtx context.Context, userId, groupId int64) error {
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
	gtx context.Context, userId, serviceId int64) ([]string, error) {
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

func (pgu *UserStorage) Exists(
	gtx context.Context, userId string) (bool, error) {
	return pgu.gd.Exists(gtx, "idx_user", "user_id", userId)
}

func (pgu *UserStorage) Count(
	gtx context.Context, filter *data.Filter) (int64, error) {
	return pgu.gd.Count(gtx, "idx_user", filter)
}

// func (us *UserStorage) SetPassword(userId, password string) error {
// 	// TODO - implement
// 	return nil
// }

// func (us *UserStorage) UpdatePassword(userId, oldPw, newPw string) error {
// 	// TODO - implement
// 	return nil
// }

// func (us *UserStorage) Verify(userId, password string) error {
// 	// TODO - implement
// 	return nil
// }
