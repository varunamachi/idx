package controller

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
)

// TODO - implement
type Service struct {
	srvStore core.ServiceStorage
}

func NewServiceController(ss core.ServiceStorage) *Service {
	return &Service{
		srvStore: ss,
	}
}

func (sc Service) Storage() core.ServiceStorage {
	return sc.srvStore
}

func (sc Service) Save(gtx context.Context, service *core.Service) error {

	return nil
}

func (sc Service) Update(gtx context.Context, service *core.Service) error {
	return nil
}

func (sc Service) GetOne(
	gtx context.Context, id int) (*core.Service, error) {
	return nil, nil
}

func (sc Service) Remove(gtx context.Context, id int) error {
	return nil
}

func (sc Service) Get(
	gtx context.Context, params *data.CommonParams) ([]*core.Service, error) {
	return nil, nil
}

func (sc Service) Exists(gtx context.Context, id int) (bool, error) {
	return false, nil
}

func (sc Service) Count(
	gtx context.Context, filter *data.Filter) (int64, error) {
	return 0, nil
}
