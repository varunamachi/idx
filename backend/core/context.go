package core

import (
	"context"

	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/event"
	"github.com/varunamachi/libx/email"
	"github.com/varunamachi/libx/httpx"
)

type Services struct {
	EventService event.Service
	MailProvider email.Provider
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

func EventService(gtx context.Context) event.Service {
	return srvs(gtx).EventService
}

func NewEventAdder(gtx context.Context, op string, data data.M) *event.Adder {

	user := httpx.GetUser[auth.User](gtx)
	userId := int64(-1)
	if user != nil {
		userId = user.Id()
	}

	return event.NewAdder(
		gtx, EventService(gtx), op, userId, data)
}

func CopyServices(source, target context.Context) context.Context {
	s := srvs(source)
	return context.WithValue(target, servicesKey, s)
}
