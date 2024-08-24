package grpdx

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
)

type groupCtl struct {
	gstore   core.GroupStorage
	svcStore core.ServiceStorage
}

func NewGroupController(
	gstore core.GroupStorage,
	svcStore core.ServiceStorage) core.GroupController {
	return &groupCtl{
		gstore:   gstore,
		svcStore: svcStore,
	}
}

func (gc *groupCtl) Storage() core.GroupStorage {
	return gc.gstore
}

func (gc *groupCtl) Save(
	gtx context.Context,
	group *core.Group) (int64, error) {

	ev := core.NewEventAdder(gtx, "group.save", data.M{
		"group": group,
	})
	id, err := gc.gstore.Save(gtx, group)
	if err != nil {
		return id, ev.Commit(err)
	}

	// if err := gc.gstore.SetPermissions(gtx, group.Id, perms); err != nil {
	// 	return ev.Commit(err)
	// }

	return id, ev.Commit(nil)
}

func (gc *groupCtl) SaveWithPerms(
	gtx context.Context,
	group *core.Group,
	perms []string) (int64, error) {

	ev := core.NewEventAdder(gtx, "group.save", data.M{
		"group": group,
	})

	id, err := gc.gstore.Save(gtx, group)
	if err != nil {
		return id, ev.Commit(err)
	}

	if err := gc.gstore.SetPermissions(gtx, group.Id, perms); err != nil {
		return id, ev.Commit(err)
	}

	return id, ev.Commit(nil)
}

func (gc *groupCtl) Update(gtx context.Context, group *core.Group) error {
	ev := core.NewEventAdder(gtx, "group.update", data.M{
		"group": group,
	})
	if err := gc.gstore.Update(gtx, group); err != nil {
		return ev.Commit(err)
	}
	return ev.Commit(nil)
}

func (gc *groupCtl) GetOne(gtx context.Context, id int64) (*core.Group, error) {
	group, err := gc.gstore.GetOne(gtx, id)
	if err != nil {
		return nil, core.NewEventAdder(gtx, "group.getOne", data.M{
			"groupId": id,
		}).Commit(err)
	}
	return group, nil
}

func (gc *groupCtl) Remove(gtx context.Context, id int64) error {
	ev := core.NewEventAdder(gtx, "group.remove", data.M{
		"groupId": id,
	})
	if err := gc.gstore.Remove(gtx, id); err != nil {
		return ev.Commit(err)
	}
	return ev.Commit(nil)
}

func (gc *groupCtl) Get(
	gtx context.Context,
	params *data.CommonParams) ([]*core.Group, error) {

	group, err := gc.gstore.Get(gtx, params)
	if err != nil {
		return nil, core.NewEventAdder(gtx, "group.get", data.M{
			"params": params,
		}).Commit(err)
	}
	return group, nil
}

func (gc *groupCtl) Exists(gtx context.Context, id int64) (bool, error) {
	group, err := gc.gstore.Exists(gtx, id)
	if err != nil {
		return false, core.NewEventAdder(gtx, "group.exists", data.M{
			"groupId": id,
		}).Commit(err)
	}
	return group, nil
}

func (gc *groupCtl) Count(
	gtx context.Context, filter *data.Filter) (int64, error) {
	group, err := gc.gstore.Count(gtx, filter)
	if err != nil {
		return 0, core.NewEventAdder(gtx, "group.count", data.M{
			"filter": filter,
		}).Commit(err)
	}
	return group, nil
}

func (gc *groupCtl) SetPermissions(
	gtx context.Context,
	groupId int64,
	perms []string) error {
	ev := core.NewEventAdder(gtx, "group.setPerms", data.M{
		"groupId": groupId,
		"perms":   perms,
	})

	if err := gc.gstore.SetPermissions(gtx, groupId, perms); err != nil {
		return ev.Commit(err)
	}

	return ev.Commit(nil)
}

func (gc *groupCtl) GetPermissions(
	gtx context.Context,
	groupId int64) ([]string, error) {
	perms, err := gc.gstore.GetPermissions(gtx, groupId)
	if err != nil {
		return nil, core.NewEventAdder(gtx, "group.count", data.M{
			"groupId": groupId,
		}).Commit(err)
	}
	return perms, nil
}

func (gc *groupCtl) AddToGroups(
	gtx context.Context, userId int64, groupId ...int64) error {
	err := gc.gstore.AddToGroups(gtx, userId, groupId...)
	return core.NewEventAdder(gtx, "user.addToGroup", data.M{
		"userId":  userId,
		"groupId": groupId,
	}).Commit(err)
}

func (gc *groupCtl) RemoveFromGroup(
	gtx context.Context, userId, groupId int64) error {
	err := gc.gstore.RemoveFromGroup(gtx, userId, groupId)
	return core.NewEventAdder(gtx, "user.removeFromGroup", data.M{
		"userId":  userId,
		"groupId": groupId,
	}).Commit(err)
}
