package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/auth"
	idxAuth "github.com/varunamachi/idx/auth"
	"github.com/varunamachi/idx/cmd"
	"github.com/varunamachi/idx/controller"
	"github.com/varunamachi/idx/core"
	idxpg "github.com/varunamachi/idx/pg"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/email"
	"github.com/varunamachi/libx/rt"
)

func main() {

	gtx, cancel := rt.Gtx()
	defer cancel()

	emailProvider, err := email.NewProviderFromEnv("IDX")
	if err != nil {
		log.Fatal().Err(err).Msg("failed initilize email provider")
	}
	evtSrv := idxpg.NewEventService()

	gd := pg.NewGetterDeleter()
	userStore := idxpg.NewUserStorage(gd)
	groupStore := idxpg.NewGroupStorage(gd)
	serviceStore := idxpg.NewServiceStorage(gd)

	hasher := auth.NewArgon2Hasher()
	credStorage := idxpg.NewCredentialStorage(hasher)

	authr := idxAuth.NewAuthenticator(credStorage)
	uctlr := controller.NewUserController(userStore, credStorage, emailProvider)
	gctlr := controller.NewGroupController(groupStore, serviceStore)
	sctlr := controller.NewServiceController(
		serviceStore, userStore, groupStore)

	gtx = core.NewContext(gtx, &core.Services{
		UserCtlr:      uctlr,
		ServiceCtlr:   sctlr,
		GroupCtlr:     gctlr,
		Authenticator: authr,
		MailProvider:  emailProvider,
		EventService:  evtSrv,
	})

	app := libx.NewApp(
		"idx", "Simple Identity Service", "0.0.1", "varunamachi").
		WithCommands(cmd.ServeCommand())

	if err := app.RunContext(gtx, os.Args); err != nil {
		log.Fatal().Err(err).Msg("exiting...")
	}

}
