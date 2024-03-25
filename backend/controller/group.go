package controller

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
)

type Group struct {
	gstore   core.GroupStorage
	svcStore core.ServiceStorage
}

func NewGroupController(
	gstore core.GroupStorage,
	svcStore core.ServiceStorage) *Group {
	return &Group{
		gstore:   gstore,
		svcStore: svcStore,
	}
}

func (gc *Group) Storage() core.GroupStorage {
	return gc.gstore
}

func (gc *Group) Save(gtx context.Context, group *core.Group) error {
	return nil
	// ev := core.NewEventAdder(gtx, "group.save", data.M{
	// 	"group": group,
	// })
	// gc.gstore.Save()
}

func (gc *Group) Update(gtx context.Context, group *core.Group) error {
	return nil
}

func (gc *Group) GetOne(gtx context.Context, id int64) (*core.Group, error) {
	return nil, nil
}

func (gc *Group) Remove(gtx context.Context, id int64) error {
	return nil
}

func (gc *Group) Get(
	gtx context.Context, params *data.CommonParams) ([]*core.Group, error) {
	return nil, nil
}

func (gc *Group) Exists(gtx context.Context, id int64) (bool, error) {
	return false, nil
}

func (gc *Group) Count(
	gtx context.Context, filter *data.Filter) (int64, error) {
	return 0, nil
}
