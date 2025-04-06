package core

import (
	"context"

	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/event"
	"github.com/varunamachi/libx/email"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
)

type Services struct {
	EventService      event.Service[int64]
	MailProvider      email.Provider
	UserController    UserController
	UserAuthenticator auth.UserAuthenticator
	ServiceController ServiceController
	GroupController   GroupController
}

type serviceHolderKey string

const servicesKey = serviceHolderKey("core-services")

func NewContext(gtx context.Context, services *Services) context.Context {
	return context.WithValue(gtx, servicesKey, services)
}

func srvs(gtx context.Context) *Services {
	svs := gtx.Value(servicesKey).(*Services)
	if svs == nil {
		panic("failed get service holder from context")
	}
	return svs

}

func MailProvider(gtx context.Context) email.Provider {
	return srvs(gtx).MailProvider
}

func EventService(gtx context.Context) event.Service[int64] {
	return srvs(gtx).EventService
}

func NewEventAdder(
	gtx context.Context, op string, data data.M) *event.Adder[int64] {

	user := httpx.GetUser[auth.User](gtx)
	userId := int64(-1)
	if user != nil {
		userId = user.Id()
	}

	return event.NewAdder(
		gtx, EventService(gtx), op, userId, data)
}

func UserCtlr(gtx context.Context) UserController {
	return srvs(gtx).UserController
}

func Authenticator(gtx context.Context) auth.UserAuthenticator {
	return srvs(gtx).UserAuthenticator
}

func MustGetUser(gtx context.Context) *User {
	user := httpx.GetUser[*User](gtx)
	if user == nil {
		panic("failed get user from context")
	}
	return user
}

func GetUser(gtx context.Context) (*User, error) {
	user := httpx.GetUser[*User](gtx)
	if user == nil {
		return nil, errx.Fmt("failed get user from context")
	}
	return user, nil
}

func ServiceCtlr(gtx context.Context) ServiceController {
	return srvs(gtx).ServiceController
}

func GroupCtlr(gtx context.Context) GroupController {
	return srvs(gtx).GroupController
}

func CopyServices(source, target context.Context) context.Context {
	s := srvs(source)
	return context.WithValue(target, servicesKey, s)
}
