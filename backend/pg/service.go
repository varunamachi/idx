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
	id int64) (*core.Service, error) {
	var service core.Service
	err := pss.gd.GetOne(gtx, "idx_service", "id", id, &service)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (pss *ServiceStorage) Remove(
	gtx context.Context,
	id int64) error {
	if err := pss.gd.Delete(gtx, "idx_service", "id"); err != nil {
		return err
	}
	return nil
}

func (pss *ServiceStorage) Get(
	gtx context.Context,
	params *data.CommonParams) ([]*core.Service, error) {
	out := make([]*core.Service, 0, params.PageSize)

	if err := pss.gd.Get(gtx, "idx_service", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (pss *ServiceStorage) Exists(
	gtx context.Context, name string) (bool, error) {
	return pss.gd.Exists(gtx, "idx_service", "name", name)
}

func (pss *ServiceStorage) Count(
	gtx context.Context, filter *data.Filter) (int64, error) {
	return pss.gd.Count(gtx, "idx_service", filter)
}

func (pss *ServiceStorage) GetByName(
	gtx context.Context, name string) (*core.Service, error) {
	var service core.Service
	err := pss.gd.GetOne(gtx, "idx_service", "name", name, &service)
	if err != nil {
		return nil, errx.Errf(err, "failed to get service with name '%s'", name)
	}
	return &service, nil
}

func (pss *ServiceStorage) GetForOwner(
	gtx context.Context, ownerId string) ([]*core.Service, error) {
	const query = `
		SELECT * 
		FROM idx_service
		WHERE owner_id = $1
		ORDER BY updated_at DESC
	`

	services := make([]*core.Service, 0, 100)
	err := pg.Conn().SelectContext(gtx, &services, query, ownerId)
	if err != nil {
		return nil, errx.Errf(err, "failed to get services for owner '%s'")
	}
	return services, nil
}
