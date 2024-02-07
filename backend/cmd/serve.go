package cmd

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/data/pg"
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
			app := libx.MustGetApp(ctx)

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
