package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/proc"
	"github.com/varunamachi/libx/rt"
)

func main() {
	gtx, cancel := rt.Gtx()
	defer cancel()

	app := libx.NewApp(
		"idx-tester", "Simple Identity Service", "0.0.1", "varunamachi").
		WithCommands(
			runCmd(gtx),
			checkPgConnCmd(gtx),
			cleanDBCmd(),
		)

	beforeBefore := app.Before
	app.Before = func(ctx *cli.Context) error {
		if err := beforeBefore(ctx); err != nil {
			return err
		}
		style := lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true).
			Align(lipgloss.Left)
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out: proc.NewWriter("Tester", os.Stderr, style, false),
			}).With().Logger()
		return nil
	}

	if err := app.RunContext(gtx, os.Args); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			errx.PrintSomeStack(err)
			log.Fatal().Err(err).Msg("exiting...")
		}
	}
}
