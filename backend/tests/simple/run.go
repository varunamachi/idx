package simple

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/client"
	"github.com/varunamachi/libx/netx"
)

func Run(gtx context.Context) error {

	err := netx.WaitForPorts(gtx, "localhost:8888", 10*time.Second)
	if err != nil {
		return err
	}

	cnt := client.New("http://localhost:8888").WithTimeout(5 * time.Minute)
	_, err = cnt.Register(gtx, super.user, super.password)
	if err != nil {
		return err
	}

	user, err := cnt.Login(gtx, super.user.UserId, super.password)
	if err != nil {
		return err
	}
	log.Info().Str("userId", user.UserId).Msg("logged in")

	log.Info().Msg("simple test successful")

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
