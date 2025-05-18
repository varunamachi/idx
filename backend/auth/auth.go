package auth

import (
	"context"
	"errors"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/errx"
)

// const (
// 	TypeUser    = "user"
// 	TypeService = "service"
// )

// type cred struct {
// 	Id       string
// 	Password string
// }

type authenticator struct {
	us core.UserController
	cs core.SecretStorage
}

func NewAuthenticator(
	us core.UserController, cs core.SecretStorage) auth.UserAuthenticator {
	return &authenticator{
		cs: cs,
		us: us,
	}
}

// Authenticate implements auth.Authenticator.
func (athn *authenticator) Authenticate(
	gtx context.Context, authData auth.AuthData) error {
	var creds core.Creds
	if err := authData.Decode(&creds); err != nil {
		return errx.Wrap(err)
	}

	if err := athn.cs.Authenticate(gtx, &creds); err != nil {
		return errx.Wrap(err)
	}
	return nil
}

func (athn *authenticator) GetUser(
	gtx context.Context, authData auth.AuthData) (auth.User, error) {

	var creds core.Creds
	if err := authData.Decode(&creds); err != nil {
		return nil, err
	}
	if creds.Type != core.AuthUser {
		return nil, errx.Errf(errors.New("invalid user for auth"),
			"entity '%s' of type '%s' cannot be authenticated as a user")
	}
	return athn.us.ByUsername(gtx, creds.UniqueName)
}
