package controller

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
)

type UserController struct {
	storage core.UserStorage
}

func (uc *UserController) Register(
	gtx context.Context, user *core.User, password string) error {
	return nil
}

func (uc *UserController) Verify(
	gtx context.Context, userId, verToken string) error {
	return nil
}

func (uc *UserController) Update(gtx context.Context, user *core.User) error {
	err := uc.storage.Update(gtx, user)
	return err
}

func (uc *UserController) GetOne(
	gtx context.Context, id int) (*core.User, error) {
	out, err := uc.storage.GetOne(gtx, id)
	if err != nil {
		core.NewEventAdder(gtx, "user.getOne", data.M{"userId": id}).
			Commit(err)

	}
	return out, err
}

func (uc *UserController) GetByUserId(
	gtx context.Context, id string) (*core.User, error) {
	out, err := uc.storage.GetByUserId(gtx, id)
	if err != nil {
		core.NewEventAdder(gtx, "user.getByUserId", data.M{"userId": id}).
			Commit(err)

	}
	return out, err
}

func (uc *UserController) SetState(
	gtx context.Context, id int, state core.UserState) error {
	err := uc.storage.SetState(gtx, id, state)
	return core.NewEventAdder(gtx, "user.setState", data.M{
		"userId": id,
		"state":  state,
	}).Commit(err)

}

func (uc *UserController) Remove(gtx context.Context, id int) error {
	err := uc.storage.Remove(gtx, id)
	return core.NewEventAdder(gtx, "user.delete", data.M{
		"userId": id,
	}).Commit(err)
}

func (uc *UserController) Get(
	gtx context.Context, params *data.CommonParams) ([]*core.User, error) {
	out, err := uc.storage.Get(gtx, params)
	if err != nil {
		core.NewEventAdder(gtx, "user.getAll", data.M{"filter": params.Filter}).
			Commit(err)
	}
	return out, err
}

func (uc *UserController) AddToGroup(
	gtx context.Context, userId, groupId int) error {
	err := uc.storage.AddToGroup(gtx, userId, groupId)
	return core.NewEventAdder(gtx, "user.addToGroup", data.M{
		"userId":  userId,
		"groupId": groupId,
	}).Commit(err)
}

func (uc *UserController) RemoveFromGroup(
	gtx context.Context, userId, groupId int) error {
	err := uc.storage.RemoveFromGroup(gtx, userId, groupId)
	return core.NewEventAdder(gtx, "user.removeFromGroup", data.M{
		"userId":  userId,
		"groupId": groupId,
	}).Commit(err)
}

func (uc *UserController) GetPermissionForService(
	gtx context.Context, userId, serviceId int) ([]string, error) {
	perms, err := uc.storage.GetPermissionForService(gtx, userId, serviceId)
	if err != nil {
		core.NewEventAdder(gtx, "user.getPerms", data.M{
			"userId":    userId,
			"serviceId": serviceId,
		}).Commit(err)
	}
	return perms, err
}
