package main

import (
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/idx/tests"
	"github.com/varunamachi/idx/tests/simple"
	"github.com/varunamachi/libx/data/pg"
)

func testSimpleCmd() *cli.Command {
	return pg.Wrap(&cli.Command{
		Name:        "simple",
		Description: "Run simple tests",
		Usage:       "Run simple tests",
		Action: func(ctx *cli.Context) error {

			if err := tests.Setup(ctx.Context); err != nil {
				return nil
			}

			defer func() {
				if err := tests.Destroy(ctx.Context); err != nil {
					log.Fatal().Err(err).Msg("failed to destroy test setup")
				}
			}()

			if err := simple.Run(ctx.Context); err != nil {
				return err
			}

			return nil
		},
	})
}
