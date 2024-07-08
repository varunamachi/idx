package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/idx/tests"
	"github.com/varunamachi/idx/tests/simple"
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
				log.Fatal().Err(err)
				return err
			}
			fmt.Println("Initialized")
			return nil
		},
		After: func(ctx *cli.Context) error {
			fmt.Println("Destroying test setup")
			if err := tests.Destroy(ctx.Context); err != nil {
				log.Fatal().Err(err).Msg("failed to destroy test setup")
				return err
			}

			fmt.Println("Destroyed")
			return nil
		},
	}
}

func simpleTestCmd() *cli.Command {
	return &cli.Command{
		Name:        "simple-tests",
		Description: "Run simple test",
		Usage:       "Run simple test",
		Action: func(ctx *cli.Context) error {
			fmt.Printf("running simple test")
			return simple.Run(ctx.Context)
		},
	}
}
