package pg

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

type PgGroupStorage struct{}

func (pgs PgGroupStorage) Save(gtx context.Context, group *core.Group) error {
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

func (pgs PgGroupStorage) Update(gtx context.Context, group *core.Group) error {
	return nil
}

func (pgs PgGroupStorage) GetOne(
	gtx context.Context, id int) (*core.Group, error) {
	return nil, nil
}

func (pgs PgGroupStorage) Remove(gtx context.Context, id int) error {
	return nil
}

func (pgs PgGroupStorage) Get(params data.CommonParams) ([]*core.Group, error) {
	return nil, nil
}
