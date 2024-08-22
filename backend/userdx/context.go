package userdx

import (
	"context"

	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
)

type serviceHolderKey string

const servicesKey = serviceHolderKey("services.user")

type Services struct {
	UserCtlr UserController
	Authr    auth.UserAuthenticator
}

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

func UserCtlr(gtx context.Context) UserController {
	return srvs(gtx).UserCtlr
}

func Authenticator(gtx context.Context) auth.UserAuthenticator {
	return srvs(gtx).Authr
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

func CopyServices(source, target context.Context) context.Context {
	s := srvs(source)
	return context.WithValue(target, servicesKey, s)
}
