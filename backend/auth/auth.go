package auth

import (
	"context"

	"github.com/varunamachi/idx/model"
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
	cs model.CredentialStorage
}

func NewAuthenticator(cs model.CredentialStorage) auth.Authenticator {
	return &authenticator{
		cs: cs,
	}
}

// Authenticate implements auth.Authenticator.
func (athn *authenticator) Authenticate(
	gtx context.Context, authData auth.AuthData) error {
	var creds model.Creds
	if err := authData.Decode(&creds); err != nil {
		return err
	}

	if err := athn.cs.Verify(gtx, &creds); err != nil {
		return err
	}
	return nil
}
