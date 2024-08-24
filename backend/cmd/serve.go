package cmd

import (
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/idx/grpdx"
	"github.com/varunamachi/idx/pg/schema"
	"github.com/varunamachi/idx/svcdx"
	"github.com/varunamachi/idx/userdx"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
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
						PrintAllAccess(false).
						WithAPIs(userdx.AuthEndpoints(gtx)...).
						WithAPIs(userdx.UserEndpoints(gtx)...).
						WithAPIs(grpdx.GroupEndpoints(gtx)...).
						WithAPIs(svcdx.ServiceEndpoints(gtx)...))

			// Create schema if required
			if err := schema.Init(gtx, "test"); err != nil {
				return errx.Wrap(err)
			}

			go func() {
				<-gtx.Done()
				log.Info().Msg("stopping the server")
				app.StopServer()
			}()

			if err := app.Serve(uint32(ctx.Uint("port"))); err != nil {
				if err != http.ErrServerClosed {
					return errx.Wrap(err)
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
	ctl := core.UserCtlr(gtx)
	return ctl.ByUsername(gtx, userId)
}

func contextMiddleware(gtx context.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(etx echo.Context) error {
			newGtx := etx.Request().Context()
			newGtx = core.CopyServices(gtx, newGtx)

			etx.SetRequest(etx.Request().WithContext(newGtx))
			// fmt.Println("set")
			return next(etx)

		}
	}
}
