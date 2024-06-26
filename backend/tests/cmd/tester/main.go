package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/tests"
	"github.com/varunamachi/libx"
)

func main() {
	gtx, cancel := tests.Gtx()
	defer cancel()

	app := libx.NewApp(
		"idx", "Simple Identity Service", "0.0.1", "varunamachi").
		WithCommands(testSimpleCmd())

	if err := app.RunContext(gtx, os.Args); err != nil {
		log.Fatal().Err(err).Msg("exiting...")
	}
}
