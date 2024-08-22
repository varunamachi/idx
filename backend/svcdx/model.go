package svcdx

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/idx/userdx"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
)

type Service struct {
	core.DbItem
	Name        string              `db:"name" json:"name"`
	OwnerId     int64               `db:"owner_id" json:"ownerId"`
	DisplayName string              `db:"display_name" json:"displayName"`
	Permissions auth.PermissionTree `db:"permissions" json:"permissions"`
}

type ServiceStorage interface {
	Save(gtx context.Context, service *Service) (int64, error)
	Update(gtx context.Context, service *Service) error
	GetOne(gtx context.Context, id int64) (*Service, error)
	GetByName(gtx context.Context, name string) (*Service, error)
	GetForOwner(gtx context.Context, ownerId string) ([]*Service, error)
	Remove(gtx context.Context, id int64) error
	Get(gtx context.Context, params *data.CommonParams) ([]*Service, error)

	AddAdmin(gtx context.Context, serviceId, userId int64) error
	GetAdmins(gtx context.Context, serviceId int64) ([]*userdx.User, error)
	RemoveAdmin(gtx context.Context, serviceId, userId int64) error
	IsAdmin(gtx context.Context, serviceId, userId int64) (bool, error)

	Exists(gtx context.Context, name string) (bool, error)
	Count(gtx context.Context, filter *data.Filter) (int64, error)

	GetPermissionForService(
		gtx context.Context, userId, serviceId int64) ([]string, error)
}

type ServiceController interface {
	ServiceStorage
	Storage() ServiceStorage
}
