package main

import (
	"context"
	"net/url"
	"os"
	"slices"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/idx/pg/schema"
	"github.com/varunamachi/idx/tests"
	"github.com/varunamachi/idx/tests/simple"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/proc"
)

func runCmd(gtx context.Context) *cli.Command {
	var process *os.Process
	procMan := proc.NewManager(gtx)

	return &cli.Command{
		Name:        "run",
		Description: "Run tests",
		Usage:       "Run tests",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "skip",
				Usage: "Skip a step. Valid: 'server', 'compose', 'clean'",
			},
			&cli.BoolFlag{
				Name:  "destroy-schema-after-test",
				Usage: "Destroy the schema once tests are done",
				Value: false,
			},
		},
		Subcommands: []*cli.Command{
			simpleTestCmd(),
		},
		Before: func(ctx *cli.Context) error {

			testConfig := getConfig(ctx)

			var err error
			process, err = tests.Setup(ctx.Context, procMan, &testConfig)
			if err != nil {
				// log.Error().Err(err).Msg("initialization failed")
				return errx.Errf(err, "initialization failed")
			}
			log.Info().Msg("initialized test setup")
			return nil
		},
		After: func(ctx *cli.Context) error {

			// So that necessary mails are sent
			// time.Sleep(200 * time.Millisecond)

			testConfig := getConfig(ctx)
			err := tests.Destroy(ctx.Context, procMan, &testConfig, process)
			if err != nil {
				// log.Error().Err(err).Msg("destruction failed")
				errx.Errf(err, "failed to destroy test setup")
			}
			log.Info().Msg("destroyed test setup")
			return nil
		},
	}
}

func getConfig(ctx *cli.Context) tests.TestConfig {
	ss := ctx.StringSlice("skip")

	return tests.TestConfig{
		SkipAppServer:     slices.Contains(ss, "server"),
		SkipDockerCompose: slices.Contains(ss, "compose"),
		SkipDataClean:     slices.Contains(ss, "clean"),
		DestroySchema:     ctx.Bool("destroy-schema-after-test"),
	}
}

func simpleTestCmd() *cli.Command {
	return &cli.Command{
		Name:        "simple",
		Description: "Run simple test",
		Usage:       "Run simple test",
		Action: func(ctx *cli.Context) error {
			return simple.Run(ctx.Context)
		},
	}
}

func checkPgConnCmd(gtx context.Context) *cli.Command {
	procMan := proc.NewManager(gtx)

	return &cli.Command{
		Name:        "check",
		Description: "Check something",
		Usage:       "Check something",
		Action: func(ctx *cli.Context) error {
			err := tests.RunDockerCompose(
				procMan,
				"up",
				tests.MustGetPgComposePath())
			if err != nil {
				return errx.Wrap(err)
			}

			defer func() {
				err := tests.RunDockerCompose(
					procMan,
					"down",
					tests.MustGetPgComposePath())
				if err != nil {
					log.Error().Err(err).Msg("failed to shutdown dc")
				}
			}()

			const pgUrl = "postgresql://idx:idxp@localhost:5432/idx-test?sslmode=disable"
			purl, err := url.Parse(pgUrl)
			if err != nil {
				return errx.Errf(err, "invalid pg-url in setup")
			}
			db, err := pg.Connect(ctx.Context, purl, "Asia/Kolkata")
			if err != nil {
				return errx.Wrap(err)
			}
			pg.SetDefaultConn(db)

			return pg.Conn().Ping()
		},
	}
}

func cleanDBCmd() *cli.Command {
	return &cli.Command{
		Name:        "clean-db",
		Description: "cleans data from idx's database tables",
		Usage:       "cleans data from idx's database tables",
		Action: func(ctx *cli.Context) error {
			if err := tests.ConnectToTestDB(ctx.Context); err != nil {
				return errx.Wrap(err)
			}
			if err := schema.CleanData(ctx.Context); err != nil {
				return errx.Wrap(err)
			}
			log.Info().Msg("database clean complete!")
			return nil
		},
	}
}
