package core

import (
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/sause/core"
)

type PgUserStorage struct {
}

func (pgu *PgUserStorage) Save(user *core.User) (err error) {
	return nil
}

func (pgu *PgUserStorage) Update(user *core.User) (err error) {
	return nil
}

func (pgu *PgUserStorage) GetOne(id string) (user *core.User, err error) {
	return nil, nil
}

func (pgu *PgUserStorage) SetState(id string, state core.UserState) error {
	return nil
}

func (pgu *PgUserStorage) Remove(id string) error {
	return nil
}

func (pgu *PgUserStorage) Get(params data.CommonParams) ([]*core.User, error) {
	return nil, nil
}
