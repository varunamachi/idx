package controller

import (
	"context"
	"time"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
)

// TODO - implement
type svcCtl struct {
	srvStore core.ServiceStorage
}

func NewServiceController(ss core.ServiceStorage) core.ServiceController {
	return &svcCtl{
		srvStore: ss,
	}
}

func (sc svcCtl) Storage() core.ServiceStorage {
	return sc.srvStore
}

func (sc svcCtl) Save(gtx context.Context, service *core.Service) error {
	user, err := core.GetUser(gtx)
	if err != nil {
		return err
	}
	ev := core.NewEventAdder(gtx, "service.save", data.M{"service": service})

	exists, err := sc.srvStore.Exists(gtx, service.Name)
	if err != nil {
		return err
	}
	if exists {
		return ev.Errf(core.ErrEntityExists,
			"service '%d:%s' already exists", service.Id, service.Name)
	}

	service.CreatedBy, service.UpdatedBy = user.SeqId(), user.SeqId()
	if err = sc.srvStore.Save(gtx, service); err != nil {
		return ev.Commit(err)
	}

	return ev.Commit(err)
}

func (sc svcCtl) Update(gtx context.Context, service *core.Service) error {
	ev := core.NewEventAdder(gtx, "service.update", data.M{
		"service": service,
	})
	user := core.MustGetUser(gtx)

	// TODO - make this list of admins of the service
	if service.OwnerId != user.SeqId() {
		return errx.Errf(core.ErrUnauthorized,
			"an user '%s' is not authorized to update service '%s'",
			user.UserId, service.Name,
		)
	}

	service.UpdatedBy, service.UpdatedAt = user.SeqId(), time.Now()
	err := sc.srvStore.Update(gtx, service)
	return ev.Commit(err)
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
	//TODO - implement
	return nil
}
func (sc *svcCtl) GetAdmins(
	gtx context.Context, serviceId int64) ([]*core.User, error) {
	//TODO - implement
	return nil, nil
}
func (sc *svcCtl) RemoveAdmin(
	gtx context.Context, serviceId, userId int64) error {
	//TODO - implement
	return nil
}
func (sc *svcCtl) IsAdmin(
	gtx context.Context, serviceId, userId int64) (bool, error) {
	//TODO - implement
	return false, nil
}
