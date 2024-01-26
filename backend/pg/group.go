package pg

import (
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
)

type PgGroupStorage struct{}

func (pgs PgGroupStorage) Save(user *core.Group) error {
	return nil
}

func (pgs PgGroupStorage) Update(user *core.Group) error {
	return nil
}

func (pgs PgGroupStorage) GetOne(id string) (*core.Group, error) {
	return nil, nil
}

func (pgs PgGroupStorage) Remove(id string) error {
	return nil
}

func (pgs PgGroupStorage) Get(params data.CommonParams) ([]*core.Group, error) {
	return nil, nil
}
