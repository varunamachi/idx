package schema

import (
	"context"
	"embed"

	"github.com/pressly/goose/v3"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

//go:embed migrations/*.sql
var migs embed.FS

func Init(gtx context.Context, initContext string) error {
	if err := CreateDB(gtx); err != nil {
		return err
	}
	// Any other DB initialization logic can go here
	return nil
}

func CreateDB(gtx context.Context) error {
	db := pg.Conn().DB
	goose.SetBaseFS(migs)

	if err := goose.SetDialect("postgres"); err != nil {
		return errx.Errf(err, "failed to initialize sql-migrator")
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return errx.Errf(err, "sql-migration failed")
	}

	return nil
}
