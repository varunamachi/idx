package auth

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/auth"
)

const (
	TypeUser    = "user"
	TypeService = "service"
)

type cred struct {
	Id       string
	Password string
}

type authenticator struct {
	us core.UserStorage
	cs core.SecretStorage
}

func NewAuthenticator(cs core.SecretStorage) auth.UserAuthenticator {
	return &authenticator{
		cs: cs,
	}
}

// Authenticate implements auth.Authenticator.
func (athn *authenticator) Authenticate(
	gtx context.Context, authData auth.AuthData) error {
	var creds core.Creds
	if err := authData.Decode(&creds); err != nil {
		return err
	}

	if err := athn.cs.Verify(gtx, &creds); err != nil {
		return err
	}
	return nil
}

func (athn *authenticator) GetUser(
	gtx context.Context, authData auth.AuthData) (auth.User, error) {

	userId, _, err := authData.ToUserAndPassword()
	if err != nil {
		return nil, err
	}
	return athn.us.GetByUserId(gtx, userId)
}
