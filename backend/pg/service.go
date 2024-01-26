package pg

import (
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
)

type PgServiceStorage struct{}

func (pss *PgServiceStorage) Save(service *core.Service) error {
	return nil
}

func (pss *PgServiceStorage) Update(service *core.Service) error {
	return nil
}

func (pss *PgServiceStorage) GetOne(id string) (*core.Service, error) {
	return nil, nil
}

func (pss *PgServiceStorage) Remove(id string) error {
	return nil
}

func (pss *PgServiceStorage) Get(
	params data.CommonParams) ([]*core.Service, error) {
	return nil, nil
}
