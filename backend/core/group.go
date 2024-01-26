package core

import (
	"context"

	"github.com/varunamachi/libx/data"
)

type Group struct {
	DbItem
	Name        string   `db:"name" json:"name"`
	DisplayName string   `db:"display_name" json:"displayName"`
	Description string   `db:"description" json:"description"`
	Perms       []string `json:"perms"`
}

type GroupStorage interface {
	Save(gtx context.Context, user *Group) error
	Update(gtx context.Context, user *Group) error
	GetOne(gtx context.Context, id string) (*Group, error)
	Remove(gtx context.Context, id string) error
	Get(gtx context.Context, params data.CommonParams) ([]*Group, error)
}
