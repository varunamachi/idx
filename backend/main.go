package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/idx/pg/schema"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/rt"
)

func main() {

	// gtx, cancel := context.WithCancel(context.Background())
	gtx, cancel := rt.Gtx()
	defer cancel()

	app := libx.NewApp(
		"idx", "Identity Service", "0.1.0", "varunamachi@gmail.com").
		WithBuildInfo(core.GetBuildInfo()).
		WithCommands()

	if err := schema.Init(gtx, "onServerStart"); err != nil {
		log.Fatal().Err(err).Msg("DB init failed")
	}

	if err := app.RunContext(gtx, os.Args); err != nil {
		errx.PrintSomeStack(err)
		log.Fatal().Msg("Exiting due to errors")
	}
}
