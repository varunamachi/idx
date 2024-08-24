package grpdx

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

type groupPgStorage struct {
	gd data.GetterDeleter
}

func NewGroupStorage(gd data.GetterDeleter) core.GroupStorage {
	return &groupPgStorage{
		gd: gd,
	}
}

func (pgs groupPgStorage) Save(
	gtx context.Context, group *core.Group) (int64, error) {
	query := `
		INSERT INTO idx_group (
			created_by,
			updated_by,
			service_id,
			name,
			display_name,
			description
		) VALUES (
			:created_by,
			:updated_by,
			:service_id,
			:name,
			:display_name,
			:description
		) ON CONFLICT (id) DO UPDATE SET
				created_by = EXCLUDED.created_by,
				updated_by = EXCLUDED.updated_by,
				service_id = EXCLUDED.service_id,
				name = EXCLUDED.name,
				display_name = EXCLUDED.display_name,
				description = EXCLUDED.description
		RETURNING id;
	`

	stmt, err := pg.Conn().PrepareNamed(query)
	if err != nil {
		return -1, errx.Errf(err, "failed to prepare query to save group")
	}

	var id int64
	if err = stmt.GetContext(gtx, &id, group); err != nil {
		return -1, errx.Errf(
			err, "failed to insert group '%s' to database", group.Id)
	}
	return id, nil

	// if _, err := pg.Conn().NamedExecContext(gtx, query, group); err != nil {
	// 	return -1, errx.Errf(
	// 		err, "failed to insert user '%s' to database", group.Id)
	// }
	// return 0, nil
}

func (pgs groupPgStorage) Update(gtx context.Context, group *core.Group) error {
	query := `
		UPDATE idx_group SET
			created_by = :created_by,
			updated_by = :updated_by,
			service_id = :service_id,
			name = :name,
			display_name = :display_name,
			description = :description
		WHERE id = :id
	`
	if _, err := pg.Conn().NamedExecContext(gtx, query, group); err != nil {
		return errx.Errf(
			err, "failed to insert user '%s' to database", group.Id)
	}

	return nil
}

func (pgs groupPgStorage) GetOne(
	gtx context.Context, id int64) (*core.Group, error) {
	var group core.Group
	err := pgs.gd.GetOne(gtx, "idx_group", "id", id, &group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (pgs groupPgStorage) Remove(gtx context.Context, id int64) error {
	err := pgs.gd.Delete(gtx, "idx_group", "id", id)
	if err != nil {
		return errx.Wrap(err)
	}
	return nil
}

func (pgs groupPgStorage) Get(
	gtx context.Context, params *data.CommonParams) ([]*core.Group, error) {
	groups := make([]*core.Group, 0, params.PageSize)
	err := pgs.gd.Get(gtx, "idx_group", params, &groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (pgs *groupPgStorage) Exists(
	gtx context.Context, id int64) (bool, error) {
	return pgs.gd.Exists(gtx, "idx_group", "id", id)
}

func (pgs *groupPgStorage) Count(
	gtx context.Context, filter *data.Filter) (int64, error) {
	return pgs.gd.Count(gtx, "idx_group", filter)
}

func (pgs *groupPgStorage) SetPermissions(
	gtx context.Context,
	groupId int64,
	perms []string) error {

	const query = `
		INSERT INTO group_to_perm(
			group_id,
			perm_id
		) VALUES (
			$1,
			$2
		) ON CONFLICT IGNORE
	`

	tx, err := pg.Conn().BeginTxx(gtx, &sql.TxOptions{})
	if err != nil {
		return errx.Errf(err, "failed to begin transaction")
	}
	ef := func(err error, fmtStr string, args ...any) error {
		if e := tx.Rollback(); e != nil {
			log.Error().Err(err).
				Msg("transaction rollback failed for fake users table")
		}
		return errx.Errf(err, fmtStr, args...)
	}

	for perm := range perms {
		_, err := tx.ExecContext(gtx, query)
		if err != nil {
			return ef(err, "failed to add permission '%s' to group '%s'",
				perm, groupId)
		}
	}

	if err := tx.Commit(); err != nil {
		return ef(err, "failed to commit addition of permissions to group")
	}

	return nil
}

func (pgs *groupPgStorage) GetPermissions(
	gtx context.Context,
	groupId int64) ([]string, error) {
	const query = `SELECT perm_id FROM group_to_perm WHERE group_id = $1`

	perms := make([]string, 0, 100)
	err := pg.Conn().SelectContext(gtx, &perms, query, groupId)
	if err != nil {
		return nil, errx.Errf(
			err, "failed get permissions for group '%d'", groupId)
	}
	return perms, nil
}

func (pgs *groupPgStorage) AddToGroups(
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

func (pgs *groupPgStorage) RemoveFromGroup(
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
