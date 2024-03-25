package pg

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

type GroupStorage struct {
	gd data.GetterDeleter
}

func NewGroupStorage(gd data.GetterDeleter) core.GroupStorage {
	return &GroupStorage{
		gd: gd,
	}
}

func (pgs GroupStorage) Save(gtx context.Context, group *core.Group) error {
	query := `
		INSERT INTO idx_user (
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
		) ON CONFLICT (user_id) DO UPDATE SET
				created_by = EXCLUDED.created_by,
				updated_by = EXCLUDED.updated_by,
				service_id = EXCLUDED.service_id,
				name = EXCLUDED.name,
				display_name = EXCLUDED.display_name,
				description = EXCLUDED.description
	`

	if _, err := pg.Conn().NamedExecContext(gtx, query, group); err != nil {
		return errx.Errf(
			err, "failed to insert user '%s' to database", group.Id)
	}

	return nil
}

func (pgs GroupStorage) Update(gtx context.Context, group *core.Group) error {
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

func (pgs GroupStorage) GetOne(
	gtx context.Context, id int64) (*core.Group, error) {
	var group core.Group
	err := pgs.gd.GetOne(gtx, "idx_group", "id", id, &group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (pgs GroupStorage) Remove(gtx context.Context, id int64) error {
	err := pgs.gd.Delete(gtx, "idx_group", "id", id)
	if err != nil {
		return err
	}
	return nil
}

func (pgs GroupStorage) Get(
	gtx context.Context, params *data.CommonParams) ([]*core.Group, error) {
	groups := make([]*core.Group, 0, params.PageSize)
	err := pgs.gd.Get(gtx, "idx_group", params, &groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (pgs *GroupStorage) Exists(
	gtx context.Context, id int64) (bool, error) {
	return pgs.gd.Exists(gtx, "idx_group", "id", id)
}

func (pgs *GroupStorage) Count(
	gtx context.Context, filter *data.Filter) (int64, error) {
	return pgs.gd.Count(gtx, "idx_group", filter)
}
