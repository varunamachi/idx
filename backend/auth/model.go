package auth

import (
	"context"

	"github.com/varunamachi/idx/model"
)

type Authenticator interface {
	AuthenticateUser(
		gtx context.Context,
		userId, password string) error
	AuthenticateService(gtx context.Context, serviceId, password string) error
}

type Authorizor interface {
	AuthorizeUser(gtx context.Context, serviceId, userId string) (*model.User, error)
}
