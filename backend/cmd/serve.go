package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/data/pg"
)

func serveCommand() *cli.Command {
	return pg.Wrap(&cli.Command{
		Name:        "serve",
		Usage:       "Start the idx identity service",
		Description: "Start the idx identity service",
		Action: func(ctx *cli.Context) error {

			gtx := ctx.Context
			app := libx.MustGetApp(ctx)

			go func() {
				<-gtx.Done()
				log.Info().Msg("stopping the server")
				app.StopServer()
			}()

			return nil
		},
	})
}
