package controller

import (
	"context"

	"github.com/varunamachi/idx/model"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/event"
)

type UserController struct {
	storage      model.UserStorage
	eventStorage event.Service
}

func (uc *UserController) Register(user *model.User, password string) error {
	return nil
}

func (uc *UserController) Verify(userId, verToken string) error {
	return nil
}

func (uc *UserController) Update(gtx context.Context, user *model.User) error {
	err := uc.storage.Update(gtx, user)
	return err
}

func (uc *UserController) GetOne(
	gtx context.Context, id int) (*model.User, error) {
	out, err := uc.storage.GetOne(gtx, id)
	return out, err
}

func (uc *UserController) GetByUserId(
	gtx context.Context, id string) (*model.User, error) {
	out, err := uc.storage.GetByUserId(gtx, id)
	return out, err
}

func (uc *UserController) SetState(
	gtx context.Context, id int, state model.UserState) error {
	err := uc.storage.SetState(gtx, id, state)
	return err
}

func (uc *UserController) Remove(gtx context.Context, id int) error {
	err := uc.storage.Remove(gtx, id)
	return err
}

func (uc *UserController) Get(
	gtx context.Context, params *data.CommonParams) ([]*model.User, error) {
	out, err := uc.storage.Get(gtx, params)
	return out, err
}

func (uc *UserController) AddToGroup(
	gtx context.Context, userId, groupId int) error {
	err := uc.storage.AddToGroup(gtx, userId, groupId)
	return err
}

func (uc *UserController) RemoveFromGroup(
	gtx context.Context, userId, groupId int) error {
	err := uc.storage.RemoveFromGroup(gtx, userId, groupId)
	return err
}

func (uc *UserController) GetPermissionForService(
	gtx context.Context, userId, serviceId int) ([]string, error) {
	perms, err := uc.storage.GetPermissionForService(gtx, userId, serviceId)
	return perms, err
}
