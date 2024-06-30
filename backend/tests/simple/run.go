package simple

import (
	"context"

	"github.com/varunamachi/idx/client"
)

func Run(gtx context.Context) error {

	cnt := client.New("localhost:8080")
	_, err := cnt.Register(gtx, super.user, super.password)
	if err != nil {
		return err
	}

	_, err = cnt.Login(gtx, super.user.UserId, super.password)
	if err != nil {
		return err
	}

	// TODO - make the email provider a service and receive email through
	// API, fake provider cant be just memory based

	// Create a super user
	// Register a user
	// Use fake email provider
	// Get the mail fro provider
	// Verify
	// Try to login
	// Logout
	// Change password
	// Login Again
	return nil

}
