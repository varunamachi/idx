package controller

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
	svcStore core.ServiceStorage) *groupCtl {
	return &groupCtl{
		gstore:   gstore,
		svcStore: svcStore,
	}
}

func (gc *groupCtl) Storage() core.GroupStorage {
	return gc.gstore
}

func (gc *groupCtl) Save(gtx context.Context, group *core.Group) error {
	return nil
	// ev := core.NewEventAdder(gtx, "group.save", data.M{
	// 	"group": group,
	// })
	// gc.gstore.Save()
}

func (gc *groupCtl) Update(gtx context.Context, group *core.Group) error {
	return nil
}

func (gc *groupCtl) GetOne(gtx context.Context, id int64) (*core.Group, error) {
	return nil, nil
}

func (gc *groupCtl) Remove(gtx context.Context, id int64) error {
	return nil
}

func (gc *groupCtl) Get(
	gtx context.Context, params *data.CommonParams) ([]*core.Group, error) {
	return nil, nil
}

func (gc *groupCtl) Exists(gtx context.Context, id int64) (bool, error) {
	return false, nil
}

func (gc *groupCtl) Count(
	gtx context.Context, filter *data.Filter) (int64, error) {
	return 0, nil
}
