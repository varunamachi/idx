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
	UserCtlr      UserController
	ServiceCtlr   ServiceController
	GroupCtlr     GroupController
	Authenticator auth.UserAuthenticator
	EventService  event.Service
	MailProvider  email.Provider
}

type serviceHolderKey string

const servicesKey = serviceHolderKey("services")

func NewContext(gtx context.Context, services *Services) context.Context {
	return context.WithValue(gtx, servicesKey, services)
}

func GetCoreServices(gtx context.Context) *Services {
	svs := gtx.Value(servicesKey).(*Services)
	if svs == nil {
		panic("failed get service holder from context")
	}
	return svs

}

func UserCtlr(gtx context.Context) UserController {
	return GetCoreServices(gtx).UserCtlr
}

func ServiceCtlr(gtx context.Context) ServiceController {
	return GetCoreServices(gtx).ServiceCtlr
}

func GroupCtlr(gtx context.Context) GroupController {
	return GetCoreServices(gtx).GroupCtlr
}

func Authenticator(gtx context.Context) auth.UserAuthenticator {
	return GetCoreServices(gtx).Authenticator
}

func EventService(gtx context.Context) event.Service {
	return GetCoreServices(gtx).EventService
}

func MailProvider(gtx context.Context) email.Provider {
	return GetCoreServices(gtx).MailProvider
}

func NewEventAdder(gtx context.Context, op string, data data.M) *event.Adder {
	userId := "N/A"
	user, err := GetUser(gtx)
	if err == nil {
		userId = user.UserId
	}

	return event.NewAdder(
		gtx, EventService(gtx), op, userId, data)
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
