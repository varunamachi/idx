package simple

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/client"
	"github.com/varunamachi/idx/mailtmpl"
	"github.com/varunamachi/idx/tests/smsrv"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/netx"
)

func Run(gtx context.Context) error {

	// Wait for the app server and create a client
	err := netx.WaitForPorts(gtx, "localhost:8888", 10*time.Second)
	if err != nil {
		return err
	}
	cnt := client.New("http://localhost:8888").WithTimeout(5 * time.Minute)

	// Create a super user
	_, err = cnt.Register(gtx, super.user, super.password)
	if err != nil {
		return err
	}

	// Login as super super
	super, err := cnt.Login(gtx, super.user.UserId, super.password)
	if err != nil {
		return err
	}
	log.Info().Str("userId", super.UserId).Msg("logged in")

	for _, up := range users {
		uclient := client.New("http://localhost:8888").
			WithTimeout(5 * time.Minute)

		id, err := uclient.Register(gtx, up.user, up.password)
		if err != nil {
			return errx.Wrap(err)
		}
		fmt.Printf("received id: %d", id)

		mailId := up.user.EmailId + ":" +
			mailtmpl.UserAccountVerificationTemplate
		msg, err := smsrv.GetMailService().Get(
			up.user.EmailId, "to", mailId)
		if err != nil {
			return errx.Wrap(err)
		}

		urlStr, ok := msg.Data["url"].(string)
		if !ok {
			return errx.Todo("no url in mail")
		}

		err = uclient.VerifyWithUrl(gtx, urlStr)
		if err != nil {
			return errx.Wrap(err)
		}

		log.Info().Str("user", up.user.UserId).Msg("verified")

		err = cnt.Approve(gtx, up.user.UserId, up.user.Role())
		if err != nil {
			return errx.Wrap(err)
		}

		lu, err := uclient.Login(gtx, up.user.UserId, up.password)
		if err != nil {
			return errx.Wrap(err)
		}

		log.Info().Str("user", lu.UserId).Msg("logged in")

		// TODO - may be logout??
	}

	// Register a user
	// Use fake email provider
	// Get the mail fro provider
	// Verify
	// Try to login
	// Logout
	// Change password
	// Login Again

	log.Info().Msg("simple test successful")
	return nil

}
