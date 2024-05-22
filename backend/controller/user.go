package controller

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/idx/mailtmpl"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/email"
	"github.com/varunamachi/libx/errx"
)

func IsSuperUser(userId string) bool {
	// TODO - get from environment variable
	return false
}

type userCtl struct {
	ustore        core.UserStorage
	credStore     core.SecretStorage
	emailProvider email.Provider
}

func NewUserController(
	ustore core.UserStorage,
	credStore core.SecretStorage,
	emailProvider email.Provider) core.UserController {
	return &userCtl{
		ustore:        ustore,
		credStore:     credStore,
		emailProvider: emailProvider,
	}
}

func (uc *userCtl) Storage() core.UserStorage {
	return uc.ustore
}

func (uc *userCtl) CredentialStorage() core.SecretStorage {
	return uc.credStore
}

func (uc *userCtl) Register(
	gtx context.Context, user *core.User, password string) error {

	evAdder := core.NewEventAdder(gtx, "user.register", data.M{
		"user": user,
	})
	exists, err := uc.ustore.Exists(gtx, user.UserId)
	if err != nil {
		return evAdder.Commit(err)
	}
	if exists {
		err = errx.Errf(core.ErrEntityExists, "user '%s' exists", user.UserId)
		return evAdder.Commit(err)
	}

	// Enable configured super user immediately
	if IsSuperUser(user.UserId) {
		user.State = core.Active
		user.AuthzRole = auth.Super
		user.SetProp("autoApproved", true)
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

	tok := core.NewToken(user.UserId, "verfiy_account", "idx_user")
	if err := uc.credStore.StoreToken(gtx, tok); err != nil {
		err = errx.Errf(err, "failed to store user verification token")
		return evAdder.Commit(err)
	}

	// mailTemplate, err := mailtmpl.UserAccountVerificationTemplate()
	// if err != nil {
	// 	return evAdder.Commit(err)
	// }

	verificationUrl := core.ToFullUrl("user/verify", tok.Id, tok.Token)
	err = core.SendSimpleMail(
		gtx,
		user.EmailId,
		mailtmpl.UserAccountVerificationTemplate,
		data.M{
			"url": verificationUrl,
		})
	if err != nil {
		return evAdder.Commit(err)
	}

	return evAdder.Commit(nil)
}

func (uc *userCtl) Verify(
	gtx context.Context, userId, verToken string) error {
	evtAdder := core.NewEventAdder(gtx, "user.verify", data.M{
		"userId": userId,
	})

	err := uc.credStore.VerifyToken(gtx, "verify_account", userId, verToken)
	if err != nil {
		return evtAdder.Commit(err)
	}

	// Is a mail required here?
	return evtAdder.Commit(nil)
}

func (uc *userCtl) Approve(
	gtx context.Context,
	userId string,
	role auth.Role,
	groups ...int64) error {

	approver, err := core.GetUser(gtx)
	ev := core.NewEventAdder(gtx, "user.approve", data.M{
		"approver": data.Qop(approver != nil, approver.UserId, "N/A"),
		"userId":   userId,
	})
	if err != nil {
		err := errx.Errf(err, "failed to get approver information")
		ev.Commit(err)
	}

	if role == auth.Super {
		err = errx.Errf(core.ErrInvalidRole,
			"role 'Super' cannot be assigned manually")
		return ev.Commit(err)
	}

	if !auth.HasRole(approver, auth.Admin) {
		err = errx.Errf(
			core.ErrUnauthorized,
			"expect role 'admin' for approver, found '%v'",
			approver.Role())
		ev.Commit(err)
	}

	user, err := uc.ustore.GetByUserId(gtx, userId)
	if err != nil {
		return ev.Commit(err)
	}

	if user.State != core.Verfied {
		err := errx.Errf(core.ErrInvalidState,
			"only user with state 'Verified' can be approved, found %v",
			user.State)
		return ev.Commit(err)
	}

	user.State = core.Active
	user.AuthzRole = role
	if err := uc.ustore.Update(gtx, user); err != nil {
		return ev.Commit(errx.Errf(err, "failed to approve user"))
	}

	if err := uc.ustore.AddToGroups(gtx, user.SeqId(), groups...); err != nil {
		return ev.Commit(errx.Errf(err, "failed to approve user with groups"))
	}

	// mt, err := mailtmpl.UserAccountApprovedTemplate()
	// if err != nil {
	// 	return ev.Commit(err)
	// }
	err = core.SendSimpleMail(
		gtx, user.EmailId, "UserAccountApprovedTemplate",
		data.M{"loginUrl": core.ToFullUrl("/login")})
	if err != nil {
		return errx.Errf(err, "failed to notify (email) user about approval")
	}

	return ev.Commit(nil)
}

func (uc *userCtl) InitResetPassword(
	gtx context.Context, userId string) error {
	ev := core.NewEventAdder(gtx, "user.pwReset.init", data.M{
		"userId": userId,
	})

	// Get user
	user, err := uc.ustore.GetByUserId(gtx, userId)
	if err != nil {
		return ev.Commit(err)
	}

	// Make sure user is in a state that allows password reset
	if !data.OneOf(user.State, core.Active, core.Disabled) {
		return ev.Errf(core.ErrInvalidState,
			"expected user state to be one of 'Active' or 'Disabled',"+
				" found '%s'",
			user.State)
	}

	// Generate a password reset token
	tok := core.NewToken(user.UserId, "password_reset", "idx_user")
	if err := uc.credStore.StoreToken(gtx, tok); err != nil {
		return ev.Errf(err, "failed to store user password reset token")
	}

	// Get mail template
	// tmpl, err := mailtmpl.PasswordResetInitTemplate()
	// if err != nil {
	// 	return ev.Errf(err, "failed to load password reset init mail")
	// }

	// Send the verification mail
	verificationUrl := core.ToFullUrl("user/pw/reset", tok.Id, tok.Token)
	err = core.SendSimpleMail(
		gtx, user.EmailId, "PasswordResetInitTemplate",
		data.M{
			"url": verificationUrl,
		})
	if err != nil {
		return ev.Errf(err, "failed to send password reset init mail")
	}

	return nil
}

func (uc *userCtl) ResetPassword(
	gtx context.Context, userId, token, newPassword string) error {
	evtAdder := core.NewEventAdder(gtx, "user.pw.reset", data.M{
		"userId": userId,
	})

	err := uc.credStore.VerifyToken(gtx, "password_reset", userId, token)
	if err != nil {
		return evtAdder.Commit(err)
	}

	err = uc.credStore.SetPassword(gtx, &core.Creds{
		Id:       userId,
		Password: newPassword,
		Type:     "user",
	})
	if err != nil {
		return evtAdder.Commit(err)
	}

	// Is a mail required here?
	return evtAdder.Commit(nil)
}

func (uc *userCtl) UpdatePassword(gtx context.Context,
	userId, oldPassword, newPassword string) error {
	evtAdder := core.NewEventAdder(gtx, "user.pw.update", data.M{
		"userId": userId,
	})

	err := uc.credStore.Verify(gtx, &core.Creds{
		Id:       userId,
		Password: oldPassword,
		Type:     "user",
	})
	if err != nil {
		return evtAdder.Commit(err)
	}

	err = uc.credStore.SetPassword(gtx, &core.Creds{
		Id:       userId,
		Password: newPassword,
		Type:     "user",
	})
	if err != nil {
		return evtAdder.Commit(err)
	}

	// Is a mail required here?
	return evtAdder.Commit(nil)
}

func (uc *userCtl) Save(gtx context.Context, user *core.User) error {
	adr := core.NewEventAdder(gtx, "user.save", data.M{"user": user})
	err := uc.ustore.Save(gtx, user)
	return adr.Commit(err)
}

func (uc *userCtl) Update(gtx context.Context, user *core.User) error {
	adr := core.NewEventAdder(gtx, "user.update", data.M{"user": user})
	err := uc.ustore.Update(gtx, user)
	return adr.Commit(err)
}

func (uc *userCtl) GetOne(
	gtx context.Context, id int64) (*core.User, error) {
	out, err := uc.ustore.GetOne(gtx, id)
	if err != nil {
		core.NewEventAdder(gtx, "user.getOne", data.M{"userId": id}).
			Commit(err)

	}
	return out, err
}

func (uc *userCtl) GetByUserId(
	gtx context.Context, id string) (*core.User, error) {
	out, err := uc.ustore.GetByUserId(gtx, id)
	if err != nil {
		core.NewEventAdder(gtx, "user.getByUserId", data.M{"userId": id}).
			Commit(err)

	}
	return out, err
}

func (uc *userCtl) SetState(
	gtx context.Context, id int64, state core.UserState) error {
	err := uc.ustore.SetState(gtx, id, state)
	return core.NewEventAdder(gtx, "user.setState", data.M{
		"userId": id,
		"state":  state,
	}).Commit(err)

}

func (uc *userCtl) Remove(gtx context.Context, id int64) error {
	err := uc.ustore.Remove(gtx, id)
	return core.NewEventAdder(gtx, "user.delete", data.M{
		"userId": id,
	}).Commit(err)
}

func (uc *userCtl) Get(
	gtx context.Context, params *data.CommonParams) ([]*core.User, error) {
	out, err := uc.ustore.Get(gtx, params)
	if err != nil {
		core.NewEventAdder(gtx, "user.getAll", data.M{"filter": params.Filter}).
			Commit(err)
	}
	return out, err
}

func (uc *userCtl) AddToGroups(
	gtx context.Context, userId int64, groupId ...int64) error {
	err := uc.ustore.AddToGroups(gtx, userId, groupId...)
	return core.NewEventAdder(gtx, "user.addToGroup", data.M{
		"userId":  userId,
		"groupId": groupId,
	}).Commit(err)
}

func (uc *userCtl) RemoveFromGroup(
	gtx context.Context, userId, groupId int64) error {
	err := uc.ustore.RemoveFromGroup(gtx, userId, groupId)
	return core.NewEventAdder(gtx, "user.removeFromGroup", data.M{
		"userId":  userId,
		"groupId": groupId,
	}).Commit(err)
}

func (uc *userCtl) GetPermissionForService(
	gtx context.Context, userId, serviceId int64) ([]string, error) {
	perms, err := uc.ustore.GetPermissionForService(gtx, userId, serviceId)
	if err != nil {
		core.NewEventAdder(gtx, "user.getPerms", data.M{
			"userId":    userId,
			"serviceId": serviceId,
		}).Commit(err)
	}
	return perms, err
}

func (uc *userCtl) Exists(gtx context.Context, id string) (bool, error) {
	out, err := uc.ustore.Exists(gtx, id)
	if err != nil {
		core.NewEventAdder(gtx, "user.exists", data.M{"id": id}).
			Commit(err)
	}
	return out, err
}

func (uc *userCtl) Count(
	gtx context.Context, filter *data.Filter) (int64, error) {
	out, err := uc.ustore.Count(gtx, filter)
	if err != nil {
		core.NewEventAdder(gtx, "user.count", data.M{"filter": filter}).
			Commit(err)
	}
	return out, err
}
