package svcdx

import (
	"context"
	"slices"
	"time"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
)

// TODO - implement
type svcCtl struct {
	srvStore *PgServiceStorage
	// userStore core.UserStorage
}

func NewServiceController(
	ss *PgServiceStorage) core.ServiceController {
	return &svcCtl{
		srvStore: ss,
	}
}

func (sc svcCtl) Storage() *PgServiceStorage {
	return sc.srvStore
}

func (sc svcCtl) Save(
	gtx context.Context, service *core.Service) (int64, error) {
	user, err := core.GetUser(gtx)
	if err != nil {
		return -1, err
	}
	ev := core.NewEventAdder(gtx, "service.save", data.M{"service": service})

	exists, err := sc.srvStore.Exists(gtx, service.Name)
	if err != nil {
		return -1, err
	}
	if exists {
		return -1, ev.Errf(core.ErrEntityExists,
			"service '%d:%s' already exists", service.Id, service.Name)
	}

	service.CreatedBy, service.UpdatedBy = user.Id(), user.Id()

	id, err := sc.srvStore.Save(gtx, service)
	if err != nil {
		return id, ev.Commit(err)
	}

	// Make the owner admin
	if err = sc.srvStore.AddAdmin(gtx, service.Id, user.Id()); err != nil {
		return id, ev.Commit(err)
	}

	return id, ev.Commit(err)
}

func (sc svcCtl) Update(gtx context.Context, service *core.Service) error {
	ev := core.NewEventAdder(gtx, "service.update", data.M{
		"service": service,
	})
	user := core.MustGetUser(gtx)

	isAdmin, err := sc.srvStore.IsAdmin(gtx, service.Id, user.Id())
	if err != nil {
		return ev.Commit(err)
	}

	if !isAdmin {
		return ev.Errf(core.ErrUnauthorized,
			"an user '%s' is not authorized to update service '%s'",
			user.Username, service.Name,
		)
	}

	service.UpdatedBy, service.UpdatedOn = user.Id(), time.Now()
	if err := sc.srvStore.Update(gtx, service); err != nil {
		return ev.Commit(err)
	}

	return ev.Commit(nil)
}

func (sc svcCtl) GetOne(
	gtx context.Context, id int64) (*core.Service, error) {
	s, err := sc.srvStore.GetOne(gtx, id)
	if err != nil {
		return nil, core.NewEventAdder(gtx, "service.getOne", data.M{
			"id": id,
		}).Commit(err)
	}
	return s, nil
}

func (sc svcCtl) Remove(gtx context.Context, id int64) error {
	ev := core.NewEventAdder(gtx, "service.remove", data.M{
		"id": id,
	})
	err := sc.srvStore.Remove(gtx, id)
	return ev.Commit(err)
}

func (sc svcCtl) Get(
	gtx context.Context, params *data.CommonParams) ([]*core.Service, error) {
	s, err := sc.srvStore.Get(gtx, params)
	if err != nil {
		return nil, core.NewEventAdder(gtx, "service.get", data.M{
			"commonParams": params,
		}).Commit(err)
	}
	return s, nil
}

func (sc svcCtl) Exists(gtx context.Context, name string) (bool, error) {
	exists, err := sc.srvStore.Exists(gtx, name)
	if err != nil {
		return false, core.NewEventAdder(gtx, "service.exists", data.M{
			"name": name,
		}).Commit(err)
	}
	return exists, nil
}

func (sc svcCtl) Count(
	gtx context.Context, filter *data.Filter) (int64, error) {
	count, err := sc.srvStore.Count(gtx, filter)
	if err != nil {
		return 0, core.NewEventAdder(gtx, "service.count", data.M{
			"filter": filter,
		}).Commit(err)
	}
	return count, nil
}

func (sc *svcCtl) GetByName(
	gtx context.Context, name string) (*core.Service, error) {
	s, err := sc.srvStore.GetByName(gtx, name)
	if err != nil {
		return nil, core.NewEventAdder(gtx, "service.getByName", data.M{
			"name": name,
		}).Commit(err)
	}
	return s, nil
}

func (sc *svcCtl) GetForOwner(
	gtx context.Context, ownerId string) ([]*core.Service, error) {
	s, err := sc.srvStore.GetForOwner(gtx, ownerId)
	if err != nil {
		return nil, core.NewEventAdder(gtx, "service.getForOwner", data.M{
			"ownerId": ownerId,
		}).Commit(err)
	}
	return s, nil
}

func (sc *svcCtl) AddAdmin(
	gtx context.Context, serviceId, userId int64) error {
	ev := core.NewEventAdder(gtx, "service.addAdmin", data.M{
		"serviceId": serviceId,
		"adminId":   userId,
	})

	curUser, err := core.GetUser(gtx)
	if err != nil {
		return ev.Commit(err)
	}

	isAdmin, err := sc.srvStore.IsAdmin(gtx, serviceId, curUser.Id())
	if err != nil {
		return ev.Commit(err)
	}
	if !isAdmin {
		return ev.Errf(core.ErrUnauthorized,
			"user '%s' is not authorized modify service admin list")
	}

	perms, err := sc.srvStore.GetPermissionForService(gtx, userId, serviceId)
	if err != nil {
		return ev.Commit(err)
	}

	if !slices.Contains(perms, PermServiceAdmin) {
		return ev.Errf(core.ErrUnauthorized,
			"only user with 'idx.serviceAdmin' permission can be "+
				"added as admin")
	}

	err = sc.srvStore.AddAdmin(gtx, serviceId, userId)
	if err != nil {
		return ev.Commit(err)
	}
	return ev.Commit(err)
}

func (sc *svcCtl) GetAdmins(
	gtx context.Context, serviceId int64) ([]*core.User, error) {
	admins, err := sc.srvStore.GetAdmins(gtx, serviceId)
	if err != nil {
		return nil, core.NewEventAdder(gtx, "service.addAdmin", data.M{
			"serviceId": serviceId,
		}).Commit(err)
	}
	return admins, nil
}

func (sc *svcCtl) RemoveAdmin(
	gtx context.Context, serviceId, userId int64) error {
	ev := core.NewEventAdder(gtx, "service.removeAdmin", data.M{
		"serviceId": serviceId,
		"adminId":   userId,
	})

	err := sc.srvStore.RemoveAdmin(gtx, serviceId, userId)
	if err != nil {
		return ev.Commit(err)
	}

	return ev.Commit(nil)
}

func (sc *svcCtl) IsAdmin(
	gtx context.Context, serviceId, userId int64) (bool, error) {
	isAdmin, err := sc.srvStore.IsAdmin(gtx, serviceId, userId)
	if err != nil {
		return false, core.NewEventAdder(gtx, "service.isAdmin", data.M{
			"serviceId": serviceId,
			"userId":    userId,
		}).Commit(err)
	}
	return isAdmin, nil
}

func (gc *svcCtl) GetPermissionForService(
	gtx context.Context, userId, serviceId int64) ([]string, error) {
	perms, err := gc.srvStore.GetPermissionForService(gtx, userId, serviceId)
	if err != nil {
		core.NewEventAdder(gtx, "user.getPerms", data.M{
			"userId":    userId,
			"serviceId": serviceId,
		}).Commit(err)
	}
	return perms, err
}
