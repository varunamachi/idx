package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/rt"
)

func main() {
	gtx, cancel := rt.Gtx()
	defer cancel()

	app := libx.NewApp(
		"idx-tester", "Simple Identity Service", "0.0.1", "varunamachi").
		WithCommands(runCmd(), checkPgConnCmd(), cleanDBCmd())

	log.Logger = log.With().Str("app", "tester").Logger()

	if err := app.RunContext(gtx, os.Args); err != nil {
		errx.PrintSomeStack(err)
		log.Fatal().Err(err).Msg("exiting...")
	}
}
