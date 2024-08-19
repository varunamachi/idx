package simple

import (
	"context"
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
		return errx.Wrap(err)
	}
	cnt := client.New("http://localhost:8888").WithTimeout(5 * time.Minute)

	// Create a super user
	_, err = cnt.Register(gtx, super.user, super.password)
	if err != nil {
		return errx.Wrap(err)
	}

	// Login as super super
	super, err := cnt.Login(gtx, super.user.UName, super.password)
	if err != nil {
		return errx.Wrap(err)
	}
	log.Info().Str("userId", super.UName).Msg("logged in")

	for _, up := range users {
		uclient := client.New("http://localhost:8888").
			WithTimeout(5 * time.Minute)

		id, err := uclient.Register(gtx, up.user, up.password)
		if err != nil {
			return errx.Wrap(err)
		}
		up.user.DbItem.Id = id

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

		log.Info().Str("user", up.user.UName).Msg("verified")

		err = cnt.Approve(gtx, up.user.Id(), up.user.Role())
		if err != nil {
			return errx.Wrap(err)
		}

		lu, err := uclient.Login(gtx, up.user.UName, up.password)
		if err != nil {
			return errx.Wrap(err)
		}

		log.Info().Str("user", lu.UName).Msg("logged in")

		// TODO - may be logout??
	}

	log.Info().Msg("simple test successful")
	return nil

}
