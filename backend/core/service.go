package core

import (
	"context"

	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
)

type Service struct {
	DbItem
	Name        string              `db:"name" json:"name"`
	OwnerId     int                 `db:"owner_id" json:"ownerId"`
	DisplayName string              `db:"display_name" json:"displayName"`
	Permissions auth.PermissionTree `db:"permissions" json:"permissions"`
}

type ServiceStorage interface {
	Save(gtx context.Context, service *Service) error
	Update(gtx context.Context, service *Service) error
	GetOne(gtx context.Context, id int) (*Service, error)
	Remove(gtx context.Context, id int) error
	Get(gtx context.Context, params *data.CommonParams) ([]*Service, error)

	Exists(gtx context.Context, id int) (bool, error)
	Count(gtx context.Context, filter *data.Filter) (int64, error)
}

type ServiceController interface {
	ServiceStorage
	Storage() ServiceStorage
}