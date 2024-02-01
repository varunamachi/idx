package pg

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

type ServiceStorage struct {
	gd data.GetterDeleter
}

func NewServiceStorage(gd data.GetterDeleter) core.ServiceStorage {
	return &ServiceStorage{
		gd: gd,
	}
}

func (pss *ServiceStorage) Save(
	gtx context.Context, service *core.Service) error {
	query := `
		INSERT INTO idx_service (
			created_by,
			updated_by,
			name,
			owner_id,
			display_name,
			permissions	
		) VALUES (
			:created_by,
			:updated_by,
			:name,
			:owner_id,
			:display_name,
			:permissions	
		) ON CONFLICT (user_id) DO UPDATE SET
				created_by = EXCLUDED.created_by,
				updated_by = EXCLUDED.updated_by,
				name = EXCLUDED.name,
				owner_id = EXCLUDED.owner_id,
				display_name = EXCLUDED.display_name,
				permissions = EXCLUDED.permissions
	`
	if _, err := pg.Conn().NamedExecContext(gtx, query, service); err != nil {
		return errx.Errf(
			err, "failed to insert service '%s' to database", service.Id)
	}

	return nil
}

func (pss *ServiceStorage) Update(
	gtx context.Context,
	service *core.Service) error {
	query := `
		UPDATE idx_service SET
			created_by = :created_by,
			updated_by = :updated_by,
			name = :name,
			owner_id = :owner_id,
			display_name = :display_name,
			permissions = :permissions	
		WHERE id = :id	
	`
	if _, err := pg.Conn().NamedExecContext(gtx, query, service); err != nil {
		return errx.Errf(
			err, "failed to update user '%s' to database", service.Id)
	}
	return nil
}

func (pss *ServiceStorage) GetOne(
	gtx context.Context,
	id int) (*core.Service, error) {
	var service core.Service
	err := pss.gd.GetOne(gtx, "idx_service", "id", id, &service)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (pss *ServiceStorage) Remove(
	gtx context.Context,
	id int) error {
	if err := pss.gd.Delete(gtx, "idx_service", "id"); err != nil {
		return err
	}
	return nil
}

func (pss *ServiceStorage) Get(
	gtx context.Context,
	params data.CommonParams) ([]*core.Service, error) {
	out := make([]*core.Service, 0, params.PageSize)

	if err := pss.gd.Get(gtx, "idx_service", &params, &out); err != nil {
		return nil, err
	}
	return out, nil
}
