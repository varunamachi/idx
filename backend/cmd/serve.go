package cmd

import (
	"context"
	"net/http"

	idxAuth "github.com/varunamachi/idx/auth"
	idxpg "github.com/varunamachi/idx/pg"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data/pg"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/idx/restapi"
	"github.com/varunamachi/libx"
)

func ServeCommand() *cli.Command {
	return pg.Wrap(&cli.Command{
		Name:        "serve",
		Usage:       "Start the idx identity service",
		Description: "Start the idx identity service",
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:  "port",
				Value: 8888,
				Usage: "Port at which the service has to run",
			},
		},
		Action: func(ctx *cli.Context) error {

			gd := pg.NewGetterDeleter()
			userStore := idxpg.NewUserStorage(gd)
			groupStore := idxpg.NewGroupStorage(gd)
			serviceStore := idxpg.NewServiceStorage(gd)

			hasher := idxAuth.NewArgon2Hasher()
			credStorage := idxpg.NewCredentialStorage(hasher)
			authr := idxAuth.NewAuthenticator(credStorage)

			gtx := ctx.Context
			app := libx.MustGetApp(ctx).
				WithServer(8080, userGetter).
				WithEndpoints(restapi.AuthEndpoints(authr)...).
				WithEndpoints(restapi.UserEndpoints(userStore)...).
				WithEndpoints(restapi.GroupEndpoints(groupStore)...).
				WithEndpoints(restapi.ServiceEndpoints(serviceStore)...)

			go func() {
				<-gtx.Done()
				log.Info().Msg("stopping the server")
				app.StopServer()
			}()

			if err := app.Serve(uint32(ctx.Uint("port"))); err != nil {
				if err != http.ErrServerClosed {
					return err
				}
			}
			return nil
		},
	})
}

func userGetter(gtx context.Context, userId string) (auth.User, error) {
	return nil, nil
}
