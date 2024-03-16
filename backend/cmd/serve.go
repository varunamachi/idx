package cmd

import (
	"context"
	"net/http"

	"github.com/varunamachi/idx/pg/schema"
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

			gtx := ctx.Context
			schema.Init(gtx, "onServerStart")
			app := libx.MustGetApp(ctx).
				WithServer(8080, &userRetriever{}).
				WithEndpoints(restapi.AuthEndpoints(gtx)...).
				WithEndpoints(restapi.UserEndpoints(gtx)...).
				WithEndpoints(restapi.GroupEndpoints(gtx)...).
				WithEndpoints(restapi.ServiceEndpoints(gtx)...)

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

type userRetriever struct {
}

func (ug *userRetriever) GetUser(
	gtx context.Context, userId string) (auth.User, error) {
	return nil, nil
}
