package cmd

import (
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/idx/pg/schema"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/httpx"

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
			// schema.Init(gtx, "onServerStart")
			app := libx.MustGetApp(ctx).
				WithServer(
					httpx.NewServer(os.Stdout, &userRetriever{}).
						WithRootMiddlewares(contextMiddleware(gtx)).
						PrintAllAccess(true).
						WithAPIs(restapi.AuthEndpoints(gtx)...).
						WithAPIs(restapi.UserEndpoints(gtx)...).
						WithAPIs(restapi.GroupEndpoints(gtx)...).
						WithAPIs(restapi.ServiceEndpoints(gtx)...))

			// Create schema if required
			if err := schema.Init(gtx, "test"); err != nil {
				return err
			}

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

func contextMiddleware(gtx context.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(etx echo.Context) error {
			services := core.GetCoreServices(gtx)
			newGtx := etx.Request().Context()
			newGtx = core.NewContext(newGtx, services)
			etx.SetRequest(etx.Request().WithContext(newGtx))
			// fmt.Println("set")
			return next(etx)

		}
	}
}
