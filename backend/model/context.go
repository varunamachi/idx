package model

import (
	"context"

	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/event"
)

type Services struct {
	userStorage    UserStorage
	serviceStorage ServiceStorage
	groupStorage   GroupStorage
	authenticator  auth.Authenticator
	eventService   event.Service
}

type serviceHolderKey string
type userHolderKey string

const servicesKey = serviceHolderKey("services")
const userKey = userHolderKey("user")

func NewContext(gtx context.Context, services *Services) context.Context {
	return context.WithValue(gtx, servicesKey, services)
}

func services(gtx context.Context) *Services {
	svs := gtx.Value(servicesKey).(*Services)
	if svs == nil {
		panic("failed get service holder from context")
	}
	return svs

}

func UserStore(gtx context.Context) UserStorage {
	return services(gtx).userStorage
}

func ServiceStore(gtx context.Context) ServiceStorage {
	return services(gtx).serviceStorage
}

func GroupStore(gtx context.Context) GroupStorage {
	return services(gtx).groupStorage
}

func Authenticator(gtx context.Context) auth.Authenticator {
	return services(gtx).authenticator
}

func EventService(gtx context.Context) event.Service {
	return services(gtx).eventService
}

func NewEventAdder(gtx context.Context, op string, data data.M) *event.Adder {

	return event.NewAdder(gtx, EventService(gtx), op, GetUser(gtx).Id(), data)
}

func GetUser(gtx context.Context) *User {
	user := gtx.Value(userKey).(*User)
	if user == nil {
		panic("failed get user from context")
	}
	return user
}
