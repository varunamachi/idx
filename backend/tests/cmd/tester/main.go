package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/rt"
)

func main() {
	gtx, cancel := rt.Gtx()
	defer cancel()

	app := libx.NewApp(
		"idx-tester", "Simple Identity Service", "0.0.1", "varunamachi").
		WithCommands(runCmd())

	if err := app.RunContext(gtx, os.Args); err != nil {
		log.Fatal().Err(err).Msg("exiting...")
	}
}
