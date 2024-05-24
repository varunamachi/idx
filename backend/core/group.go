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
	GetOne(gtx context.Context, id int64) (*Group, error)
	Remove(gtx context.Context, id int64) error
	Get(gtx context.Context, params *data.CommonParams) ([]*Group, error)

	Exists(gtx context.Context, id int64) (bool, error)
	Count(gtx context.Context, filter *data.Filter) (int64, error)

	SetPermissions(gtx context.Context, groupId int64, perms []string) error
	GetPermissions(gtx context.Context, groupId int64) ([]string, error)

	AddToGroups(gtx context.Context, userId int64, groupIds ...int64) error
	RemoveFromGroup(gtx context.Context, userId, groupId int64) error
}

type GroupController interface {
	GroupStorage

	Storage() GroupStorage
	SaveWithPerms(gtx context.Context, group *Group, perms []string) error
}
