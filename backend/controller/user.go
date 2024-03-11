package controller

import (
	"context"

	"github.com/varunamachi/idx/model"
	"github.com/varunamachi/libx/data"
)

type UserController struct {
	storage model.UserStorage
}

func (uc *UserController) Register(user *model.User, password string) error {
	return nil
}

func (uc *UserController) Verify(userId, verToken string) error {
	return nil
}

func (uc *UserController) Update(gtx context.Context, user *model.User) error {
	return uc.storage.Update(gtx, user)
}

func (uc *UserController) GetOne(
	gtx context.Context, id int) (*model.User, error) {
	return uc.storage.GetOne(gtx, id)
}
func (uc *UserController) GetByUserId(
	gtx context.Context, id string) (*model.User, error) {
	return uc.storage.GetByUserId(gtx, id)
}
func (uc *UserController) SetState(
	gtx context.Context, id int, state model.UserState) error {
	return nil
}
func (uc *UserController) Remove(gtx context.Context, id int) error {
	return nil
}
func (uc *UserController) Get(
	gtx context.Context, params *data.CommonParams) ([]*model.User, error) {
	return nil, nil
}
func (uc *UserController) AddToGroup(
	gtx context.Context, userId, groupId int) error {
	return nil
}
func (uc *UserController) RemoveFromGroup(
	gtx context.Context, userId, groupId int) error {
	return nil
}
func (uc *UserController) GetPermissionForService(
	gtx context.Context, userId, serviceId int) ([]string, error) {
	return nil, nil
}
