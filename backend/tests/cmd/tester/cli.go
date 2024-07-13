package main

import (
	"fmt"
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/idx/tests"
	"github.com/varunamachi/idx/tests/simple"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

func runCmd() *cli.Command {
	return &cli.Command{
		Name:        "run",
		Description: "Run tests",
		Usage:       "Run tests",
		Flags:       []cli.Flag{},
		Subcommands: []*cli.Command{
			simpleTestCmd(),
		},
		Before: func(ctx *cli.Context) error {
			fmt.Println("Initializing test setup")
			if err := tests.Setup(ctx.Context); err != nil {
				// log.Error().Err(err).Msg("initialization failed")
				return errx.Errf(err, "initialization failed")
			}
			return nil
		},
		After: func(ctx *cli.Context) error {
			fmt.Println("Destroying test setup")
			if err := tests.Destroy(ctx.Context); err != nil {
				// log.Error().Err(err).Msg("destruction failed")
				errx.Errf(err, "failed to destroy test setup")
			}
			return nil
		},
	}
}

func simpleTestCmd() *cli.Command {
	return &cli.Command{
		Name:        "simple",
		Description: "Run simple test",
		Usage:       "Run simple test",
		Action: func(ctx *cli.Context) error {
			fmt.Printf("running simple test")
			return simple.Run(ctx.Context)
		},
	}
}

func checkPgConn() *cli.Command {
	return &cli.Command{
		Name:        "check",
		Description: "Check something",
		Usage:       "Check something",
		Action: func(ctx *cli.Context) error {
			err := tests.RunDockerCompose("up",
				tests.MustGetPgDockerComposePath())
			if err != nil {
				return err
			}

			defer func() {
				err := tests.RunDockerCompose("down",
					tests.MustGetPgDockerComposePath())
				if err != nil {
					log.Error().Err(err).Msg("failed to shutdown dc")
				}
			}()

			const pgUrl = "postgresql://idx:idxp@localhost:5432/test-data?sslmode=disable"
			purl, err := url.Parse(pgUrl)
			if err != nil {
				return errx.Errf(err, "invalid pg-url in setup")
			}
			db, err := pg.Connect(ctx.Context, purl, "Asia/Kolkata")
			if err != nil {
				return err
			}
			pg.SetDefaultConn(db)

			return pg.Conn().Ping()
		},
	}
}
