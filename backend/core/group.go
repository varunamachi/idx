package core

import (
	"context"

	"github.com/varunamachi/libx/data"
)

type Group struct {
	DbItem
	ServiceId   int      `db:"service_id" json:"service_id"`
	Name        string   `db:"name" json:"name"`
	DisplayName string   `db:"display_name" json:"displayName"`
	Description string   `db:"description" json:"description"`
	Perms       []string `json:"perms"`
}

type GroupStorage interface {
	Save(gtx context.Context, group *Group) error
	Update(gtx context.Context, group *Group) error
	GetOne(gtx context.Context, id int) (*Group, error)
	Remove(gtx context.Context, id int) error
	Get(gtx context.Context, params *data.CommonParams) ([]*Group, error)

	Exists(gtx context.Context, id int) (bool, error)
	Count(gtx context.Context, filter *data.Filter) (int64, error)
}

type GroupController interface {
	GroupStorage

	Storage() GroupStorage
}
