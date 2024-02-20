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
	cs core.CredentialStorage
}

func NewAuthenticator(cs core.CredentialStorage) auth.Authenticator {
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
