package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/auth"
	idxAuth "github.com/varunamachi/idx/auth"
	"github.com/varunamachi/idx/cmd"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/idx/grpdx"
	idxpg "github.com/varunamachi/idx/pg"
	"github.com/varunamachi/idx/svcdx"
	"github.com/varunamachi/idx/userdx"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/email"
	"github.com/varunamachi/libx/rt"
)

func main() {

	// TODO - will these work if pg database is initialize later?

	gtx, cancel := rt.Gtx()
	defer cancel()

	emailProvider, err := email.ProviderFromEnv("IDX")
	if err != nil {
		log.Fatal().Err(err).Msg("failed initilize email provider")
	}
	if emailProvider == nil {
		log.Fatal().Msg("email provider not defined in environment")
	}

	evtSrv := idxpg.NewEventService()

	gd := pg.NewGetterDeleter()
	userStore := userdx.NewUserStorage(gd)
	serviceStore := svcdx.NewServiceStorage(gd)
	groupStore := grpdx.NewGroupStorage(gd)

	hasher := auth.NewArgon2Hasher()
	credStorage := userdx.NewCredentialStorage(hasher)

	authr := idxAuth.NewAuthenticator(userStore, credStorage)
	uctlr := userdx.NewUserController(userStore, credStorage, emailProvider)
	sctlr := svcdx.NewServiceController(serviceStore, userStore)
	gctlr := grpdx.NewGroupController(groupStore, serviceStore)

	gtx = core.NewContext(gtx, &core.Services{
		UserController:    uctlr,
		ServiceController: sctlr,
		GroupController:   gctlr,
		UserAuthenticator: authr,
		MailProvider:      emailProvider,
		EventService:      evtSrv,
	})

	app := libx.NewApp(
		"idx", "Simple Identity Service", "0.0.1", "varunamachi").
		WithCommands(cmd.ServeCommand())

	if err := app.RunContext(gtx, os.Args); err != nil {
		log.Fatal().Err(err).Msg("exiting...")
	}

}
