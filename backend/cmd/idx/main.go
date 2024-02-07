package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/cmd"
	idxpg "github.com/varunamachi/idx/pg"
	"github.com/varunamachi/idx/restapi"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/rt"
)

func main() {

	gtx, cancel := rt.Gtx()
	defer cancel()

	gd := pg.NewGetterDeleter()
	userStore := idxpg.NewUserStorage(gd)
	groupStore := idxpg.NewGroupStorage(gd)
	serviceStore := idxpg.NewServiceStorage(gd)

	app := libx.NewApp(
		"idx", "Simple Identity Service", "0.0.1", "varunamachi")
	app.WithServer(8080, userGetter).
		WithEndpoints(restapi.UserEndpoints(userStore)...).
		WithEndpoints(restapi.GroupEndpoints(groupStore)...).
		WithEndpoints(restapi.ServiceEndpoints(serviceStore)...).
		WithCommands(cmd.ServeCommand())

	if err := app.RunContext(gtx, os.Args); err != nil {
		log.Fatal().Err(err).Msg("exiting...")
	}

}

func userGetter(gtx context.Context, userId string) (auth.User, error) {
	return nil, nil
}
