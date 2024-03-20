package controller

import (
	"context"
	"errors"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/email"
	"github.com/varunamachi/libx/errx"
)

func IsSuperUser(userId string) bool {
	// TODO - get from environment variable
	return false
}

var (
	ErrUserExists = errors.New("user.exists")
)

type User struct {
	ustore        core.UserStorage
	credStore     core.SecretStorage
	emailProvider email.Provider
}

func NewUserController(
	ustore core.UserStorage,
	credStore core.SecretStorage,
	emailProvider email.Provider) *User {
	return &User{
		ustore:        ustore,
		credStore:     credStore,
		emailProvider: emailProvider,
	}
}

func (uc *User) Storage() core.UserStorage {
	return uc.ustore
}

func (uc *User) CredentialStorage() core.SecretStorage {
	return uc.credStore
}

func (uc *User) Register(
	gtx context.Context, user *core.User, password string) error {

	evAdder := core.NewEventAdder(gtx, "user.register", data.M{
		"user": user,
	})
	exists, err := uc.ustore.Exists(gtx, user.UserId)
	if err != nil {
		return evAdder.Commit(err)
	}
	if exists {
		err = errx.Errf(ErrUserExists, "user '%s' exists", user.UserId)
		evAdder.Commit(err)
	}

	// Enable configured super user immediately
	if IsSuperUser(user.UserId) {
		user.State = core.Active
		user.AuthzRole = auth.Super
	} else {
		user.State = core.Created
		user.AuthzRole = auth.Normal
	}

	if err := uc.ustore.Save(gtx, user); err != nil {
		return evAdder.Commit(err)
	}

	creds := &core.Creds{
		Id:       user.UserId,
		Password: password,
		Type:     "user",
	}
	if err := uc.credStore.SetPassword(gtx, creds); err != nil {
		return evAdder.Commit(err)
	}

	// tok := core.NewToken()
	// TODO -
	// - Generate token
	// - Store token
	// - Generate link to verify the token
	// - Send the link in the email

	return evAdder.Commit(nil)
}

func (uc *User) Verify(
	gtx context.Context, userId, verToken string) error {
	return nil
}

func (uc *User) Approve(gtx context.Context, userId string) error {
	return nil
}

func (uc *User) InitResetPassword(
	gtx context.Context, userId string) error {
	// TODO - implement
	return nil
}

func (uc *User) ResetPassword(
	gtx context.Context, userId, token, newPassword string) error {
	// TODO - implement
	return nil
}

func (uc *User) UpdatePassword(gtx context.Context,
	userId, oldPassword, newPassword string) error {
	// TODO - implement
	return nil
}

func (uc *User) Save(gtx context.Context, user *core.User) error {
	adr := core.NewEventAdder(gtx, "user.save", data.M{"user": user})
	err := uc.ustore.Save(gtx, user)
	return adr.Commit(err)
}

func (uc *User) Update(gtx context.Context, user *core.User) error {
	adr := core.NewEventAdder(gtx, "user.update", data.M{"user": user})
	err := uc.ustore.Update(gtx, user)
	return adr.Commit(err)
}

func (uc *User) GetOne(
	gtx context.Context, id int) (*core.User, error) {
	out, err := uc.ustore.GetOne(gtx, id)
	if err != nil {
		core.NewEventAdder(gtx, "user.getOne", data.M{"userId": id}).
			Commit(err)

	}
	return out, err
}

func (uc *User) GetByUserId(
	gtx context.Context, id string) (*core.User, error) {
	out, err := uc.ustore.GetByUserId(gtx, id)
	if err != nil {
		core.NewEventAdder(gtx, "user.getByUserId", data.M{"userId": id}).
			Commit(err)

	}
	return out, err
}

func (uc *User) SetState(
	gtx context.Context, id int, state core.UserState) error {
	err := uc.ustore.SetState(gtx, id, state)
	return core.NewEventAdder(gtx, "user.setState", data.M{
		"userId": id,
		"state":  state,
	}).Commit(err)

}

func (uc *User) Remove(gtx context.Context, id int) error {
	err := uc.ustore.Remove(gtx, id)
	return core.NewEventAdder(gtx, "user.delete", data.M{
		"userId": id,
	}).Commit(err)
}

func (uc *User) Get(
	gtx context.Context, params *data.CommonParams) ([]*core.User, error) {
	out, err := uc.ustore.Get(gtx, params)
	if err != nil {
		core.NewEventAdder(gtx, "user.getAll", data.M{"filter": params.Filter}).
			Commit(err)
	}
	return out, err
}

func (uc *User) AddToGroup(
	gtx context.Context, userId, groupId int) error {
	err := uc.ustore.AddToGroup(gtx, userId, groupId)
	return core.NewEventAdder(gtx, "user.addToGroup", data.M{
		"userId":  userId,
		"groupId": groupId,
	}).Commit(err)
}

func (uc *User) RemoveFromGroup(
	gtx context.Context, userId, groupId int) error {
	err := uc.ustore.RemoveFromGroup(gtx, userId, groupId)
	return core.NewEventAdder(gtx, "user.removeFromGroup", data.M{
		"userId":  userId,
		"groupId": groupId,
	}).Commit(err)
}

func (uc *User) GetPermissionForService(
	gtx context.Context, userId, serviceId int) ([]string, error) {
	perms, err := uc.ustore.GetPermissionForService(gtx, userId, serviceId)
	if err != nil {
		core.NewEventAdder(gtx, "user.getPerms", data.M{
			"userId":    userId,
			"serviceId": serviceId,
		}).Commit(err)
	}
	return perms, err
}

func (uc *User) Exists(gtx context.Context, id string) (bool, error) {
	out, err := uc.ustore.Exists(gtx, id)
	if err != nil {
		core.NewEventAdder(gtx, "user.exists", data.M{"id": id}).
			Commit(err)
	}
	return out, err
}

func (uc *User) Count(gtx context.Context, filter *data.Filter) (int64, error) {
	out, err := uc.ustore.Count(gtx, filter)
	if err != nil {
		core.NewEventAdder(gtx, "user.count", data.M{"filter": filter}).
			Commit(err)
	}
	return out, err
}
