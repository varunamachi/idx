package schema

import "github.com/urfave/cli/v2"

func Commands() []*cli.Command {
	return []*cli.Command{
		DBInitCmd(),
		DBDestroyCmd(),
	}
}

func DBInitCmd() *cli.Command {
	return &cli.Command{
		Name:        "db-init",
		Description: "Initialize database schema etc.",
		Usage:       "Initialize database schema etc.",
		Action: func(ctx *cli.Context) error {
			return Init(ctx.Context, "explicit")
		},
	}
}

func DBDestroyCmd() *cli.Command {
	return &cli.Command{
		Name:        "db-destroy",
		Description: "Deletes the schema",
		Usage:       "Deletes the schema",
	}
}
