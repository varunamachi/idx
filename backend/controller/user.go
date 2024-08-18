package controller

import (
	"context"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/idx/mailtmpl"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/email"
	"github.com/varunamachi/libx/errx"
)

var roleMapping = map[string]auth.Role{}

func mappedRole(userId string) auth.Role {
	roleEnv := os.Getenv("IDX_ROLE_MAPPING")
	if roleEnv == "" {
		return auth.None
	}

	if len(roleMapping) == 0 {
		fields := strings.Split(roleEnv, ",")
		for _, f := range fields {

			asocs := strings.Split(f, ":")
			if len(asocs) < 2 {
				log.Error().Str("field", f).Msg("invalid role association")
				continue
			}
			u := asocs[0]
			r := asocs[1]

			role := auth.ToRole(r)
			if role == auth.None {
				log.Error().Str("roleStr", r).
					Msg("invalid role association")
				continue
			}

			roleMapping[u] = role
		}

	}

	if role, found := roleMapping[userId]; found {
		return role
	}
	return auth.None
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
	gtx context.Context, user *core.User, password string) (int64, error) {

	evAdder := core.NewEventAdder(gtx, "user.register", data.M{
		"user": user,
	})
	exists, err := uc.ustore.Exists(gtx, user.UName)
	if err != nil {
		return -1, evAdder.Commit(err)
	}
	if exists {
		err = errx.Errf(core.ErrEntityExists, "user '%s' exists", user.UName)
		return -1, evAdder.Commit(err)
	}

	// Enable configured super user immediately
	autoApproved := false
	if role := mappedRole(user.UName); role != auth.None {
		user.State = core.Active
		user.AuthzRole = role
		autoApproved = true
		user.SetProp("autoApproved", true)
	} else {
		user.State = core.Created
		user.AuthzRole = auth.Normal
	}

	id, err := uc.ustore.Save(gtx, user)
	if err != nil {
		return id, evAdder.Commit(err)
	}

	creds := &core.Creds{
		UniqueName: user.UName,
		Password:   password,
		Type:       "user",
	}
	if err := uc.credStore.SetPassword(gtx, creds); err != nil {
		return id, evAdder.Commit(err)
	}

	// No need to verify auto approved accounts
	if user.State == core.Active {
		return id, nil
	}

	if !autoApproved {
		tok := core.NewToken(user.UName, "verfiy_account", "idx_user")
		if err := uc.credStore.StoreToken(gtx, tok); err != nil {
			err = errx.Errf(err, "failed to store user verification token")
			return id, evAdder.Commit(err)
		}

		verificationUrl := core.ToFullUrl(
			tok.UniqueName, "api/v1/user/verify", tok.Token)
		err = core.SendSimpleMail(
			gtx,
			user.EmailId,
			mailtmpl.UserAccountVerificationTemplate,
			data.M{
				"url": verificationUrl,
			})
		if err != nil {
			return id, evAdder.Commit(err)
		}
	}

	return id, evAdder.Commit(nil)
}

func (uc *userCtl) Verify(
	gtx context.Context, userName, verToken string) error {
	evtAdder := core.NewEventAdder(gtx, "user.verify", data.M{
		"userId": userName,
	})

	err := uc.credStore.VerifyToken(gtx, userName, "verify_account", verToken)
	if err != nil {
		return evtAdder.Commit(err)
	}

	// set the state
	// err := uc.ustore.SetState(gtx)

	// Is a mail required here?
	return evtAdder.Commit(nil)
}

func (uc *userCtl) Approve(
	gtx context.Context,
	userId int64,
	role auth.Role,
	groups ...int64) error {

	approver, err := core.GetUser(gtx)
	ev := core.NewEventAdder(gtx, "user.approve", data.M{
		"approver": data.Qop(approver != nil, approver.Id(), -1),
		"userId":   userId,
	})
	if err != nil {
		err := errx.Errf(err, "failed to get approver information")
		return errx.Wrap(ev.Commit(err))
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
		return ev.Commit(err)
	}

	user, err := uc.ustore.GetOne(gtx, userId)
	if err != nil {
		return errx.Wrap(ev.Commit(err))
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

	// if err := uc.ustore.AddToGroups(gtx, user.SeqId(), groups...); err != nil {
	// 	return ev.Commit(errx.Errf(err, "failed to approve user with groups"))
	// }

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
	gtx context.Context, userName string) error {
	ev := core.NewEventAdder(gtx, "user.pwReset.init", data.M{
		"userId": userName,
	})

	// Get user
	user, err := uc.ustore.ByUsername(gtx, userName)
	if err != nil {
		return errx.Wrap(ev.Commit(err))
	}

	// Make sure user is in a state that allows password reset
	if !data.OneOf(user.State, core.Active, core.Disabled) {
		return ev.Errf(core.ErrInvalidState,
			"expected user state to be one of 'Active' or 'Disabled',"+
				" found '%s'",
			user.State)
	}

	// Generate a password reset token
	tok := core.NewToken(user.UName, "password_reset", "idx_user")
	if err := uc.credStore.StoreToken(gtx, tok); err != nil {
		return ev.Errf(err, "failed to store user password reset token")
	}

	// Get mail template
	// tmpl, err := mailtmpl.PasswordResetInitTemplate()
	// if err != nil {
	// 	return ev.Errf(err, "failed to load password reset init mail")
	// }

	// Send the verification mail
	verificationUrl := core.ToFullUrl(
		"api/v1/user/pw/reset",
		tok.UniqueName,
		tok.Token)
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
	gtx context.Context, userName, token, newPassword string) error {
	evtAdder := core.NewEventAdder(gtx, "user.pw.reset", data.M{
		"userId": userName,
	})

	err := uc.credStore.VerifyToken(gtx, userName, "password_reset", token)
	if err != nil {
		return evtAdder.Commit(err)
	}

	err = uc.credStore.SetPassword(gtx, &core.Creds{
		UniqueName: userName,
		Password:   newPassword,
		Type:       core.AuthUser,
	})
	if err != nil {
		return evtAdder.Commit(err)
	}

	// Is a mail required here?
	return evtAdder.Commit(nil)
}

func (uc *userCtl) UpdatePassword(gtx context.Context,
	userName, oldPassword, newPassword string) error {
	evtAdder := core.NewEventAdder(gtx, "user.pw.update", data.M{
		"userId": userName,
	})

	err := uc.credStore.Verify(gtx, &core.Creds{
		UniqueName: userName,
		Password:   oldPassword,
		Type:       core.AuthUser,
	})
	if err != nil {
		return evtAdder.Commit(err)
	}

	err = uc.credStore.SetPassword(gtx, &core.Creds{
		UniqueName: userName,
		Password:   newPassword,
		Type:       core.AuthUser,
	})
	if err != nil {
		return evtAdder.Commit(err)
	}

	// Is a mail required here?
	return evtAdder.Commit(nil)
}

func (uc *userCtl) Save(gtx context.Context, user *core.User) (int64, error) {
	adr := core.NewEventAdder(gtx, "user.save", data.M{"user": user})
	id, err := uc.ustore.Save(gtx, user)
	return id, adr.Commit(err)
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

func (uc *userCtl) ByUsername(
	gtx context.Context, id string) (*core.User, error) {
	out, err := uc.ustore.ByUsername(gtx, id)
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
