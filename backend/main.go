package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/rt"
	"github.com/varunamachi/marukatte-server/core"
)

func main() {

	// gtx, cancel := context.WithCancel(context.Background())
	gtx, cancel := rt.Gtx()
	defer cancel()

	bi := libx.BuildInfo{
		GitTag:    core.GitTag,
		GitHash:   core.GitHash,
		GitBranch: core.GitBranch,
		BuildTime: core.BuildTime,
		BuildHost: core.BuildHost,
		BuildUser: core.BuildUser,
	}

	app := libx.NewApp(
		"mks", "Marukatte Server", "0.3.0", "varunamachi@gmail.com").
		WithBuildInfo(&bi).
		WithCommands()

	if err := app.RunContext(gtx, os.Args); err != nil {
		errx.PrintSomeStack(err)
		log.Fatal().Msg("Exiting due to errors")
	}
}
