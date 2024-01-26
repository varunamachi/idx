package core

import (
	"context"

	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
)

type Service struct {
	DbItem
	Name        string              `db:"name" json:"name"`
	DisplayName string              `db:"display_name" json:"displayName"`
	Permissions auth.PermissionTree `db:"permissions" json:"permissions"`
}

type ServiceStorage interface {
	Save(gtx context.Context, service *Service) error
	Update(gtx context.Context, service *Service) error
	GetOne(gtx context.Context, id string) (*Service, error)
	Remove(gtx context.Context, id string) error
	Get(gtx context.Context, params data.CommonParams) ([]*Service, error)
}
