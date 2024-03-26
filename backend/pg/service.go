package pg

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

type svcPgStorage struct {
	gd data.GetterDeleter
}

func NewServiceStorage(gd data.GetterDeleter) core.ServiceStorage {
	return &svcPgStorage{
		gd: gd,
	}
}

func (pss *svcPgStorage) Save(
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

func (pss *svcPgStorage) Update(
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

func (pss *svcPgStorage) GetOne(
	gtx context.Context,
	id int64) (*core.Service, error) {
	var service core.Service
	err := pss.gd.GetOne(gtx, "idx_service", "id", id, &service)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (pss *svcPgStorage) Remove(
	gtx context.Context,
	id int64) error {
	if err := pss.gd.Delete(gtx, "idx_service", "id"); err != nil {
		return err
	}
	return nil
}

func (pss *svcPgStorage) Get(
	gtx context.Context,
	params *data.CommonParams) ([]*core.Service, error) {
	out := make([]*core.Service, 0, params.PageSize)

	if err := pss.gd.Get(gtx, "idx_service", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (pss *svcPgStorage) Exists(
	gtx context.Context, name string) (bool, error) {
	return pss.gd.Exists(gtx, "idx_service", "name", name)
}

func (pss *svcPgStorage) Count(
	gtx context.Context, filter *data.Filter) (int64, error) {
	return pss.gd.Count(gtx, "idx_service", filter)
}

func (pss *svcPgStorage) GetByName(
	gtx context.Context, name string) (*core.Service, error) {
	var service core.Service
	err := pss.gd.GetOne(gtx, "idx_service", "name", name, &service)
	if err != nil {
		return nil, errx.Errf(err, "failed to get service with name '%s'", name)
	}
	return &service, nil
}

func (pss *svcPgStorage) GetForOwner(
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

func (pss *svcPgStorage) AddAdmin(
	gtx context.Context, serviceId, userId int64) error {
	const query = `
		INSERT INTO service_to_owner(
			user_id,
			group_id
		) VALUES (
			$1,
			$2
		)
	`
	_, err := pg.Conn().ExecContext(gtx, query, serviceId, userId)
	if err != nil {
		return errx.Errf(err,
			"failed to add admin '%d' to '%d'", userId, serviceId)
	}

	return nil
}

func (pss *svcPgStorage) GetAdmins(
	gtx context.Context, serviceId int64) ([]*core.User, error) {
	const query = `
		SELECT u.* 
		FROM idx_user u
		JOIN service_to_owner so ON u.id = so.admin_id
		WHERE service_id = $1
		ORDER BY admin_id DESC
	`

	admins := make([]*core.User, 0, 100)
	err := pg.Conn().SelectContext(gtx, &admins, query, serviceId)
	if err != nil {
		return nil, errx.Errf(err, "failed to get admins for %d", serviceId)
	}

	return nil, nil
}

func (pss *svcPgStorage) RemoveAdmin(
	gtx context.Context, serviceId, userId int64) error {
	const query = `
		DELETE FROM service_to_owner WHERE service_id = $1 AND admin_id = $2
	`
	_, err := pg.Conn().ExecContext(gtx, query, serviceId, userId)
	if err != nil {
		return errx.Errf(err,
			"failed to remove admin '%d' from service '%d'",
			userId,
			serviceId,
		)
	}

	return nil
}

func (pss *svcPgStorage) IsAdmin(
	gtx context.Context, serviceId, adminId int64) (bool, error) {
	const query = `
		SELECT EXISTS( 
			SELECT 1 
			FROM service_to_owner 
			WHERE 
				service_id = $1 AND
				owner_id = $2
		)
	`

	isAdmin := false
	err := pg.Conn().SelectContext(gtx, &isAdmin, query, adminId)
	if err != nil {
		return false, errx.Errf(err,
			"failed to check if '%s' is an admin of service '%d'",
			serviceId,
			adminId,
		)
	}
	return isAdmin, nil
}
