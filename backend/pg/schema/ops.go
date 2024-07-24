package schema

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

//go:embed migrations/*.sql
var migs embed.FS

func Init(gtx context.Context, initContext string) error {
	if err := Create(gtx); err != nil {
		return err
	}
	// Any other DB initialization logic can go here
	return nil
}

func Create(gtx context.Context) error {
	db := pg.Conn().DB
	goose.SetBaseFS(migs)

	if err := goose.SetDialect("postgres"); err != nil {
		return errx.Errf(err, "failed to initialize sql-migrator")
	}

	if err := goose.UpContext(gtx, db, "migrations"); err != nil {
		return errx.Errf(err, "sql-migration failed")
	}

	return nil
}

func Destroy(gtx context.Context) error {
	if pg.Conn() == nil {
		return fmt.Errorf("db connection does not exist")
	}

	db := pg.Conn().DB
	goose.SetBaseFS(migs)

	if err := goose.SetDialect("postgres"); err != nil {
		return errx.Errf(err, "failed to initialize sql-migrator")
	}

	if err := goose.DownContext(gtx, db, "migrations"); err != nil {
		return errx.Errf(err, "sql-migration failed")
	}

	return nil
}

func CleanData(gtx context.Context) error {
	if pg.Conn() == nil {
		return errx.Errf(errors.New("no active db connection"),
			"Not connected to database")
	}

	tables := []string{
		"idx_token",
		"credential",
		"group_to_perm",
		"user_to_group",
		"idx_event",
		"user_pass",
		"service_to_owner",
		"idx_group",
		"idx_service",
		"idx_user",
	}

	for _, table := range tables {
		query, args, err := squirrel.StatementBuilder.Delete(table).ToSql()
		if err != nil {
			return errx.Errf(err,
				"failed to generate table deletion query for '%s'", table)
		}

		_, err = pg.Conn().ExecContext(gtx, query, args...)
		if err != nil {
			return errx.Errf(err, "failed to clean data from table '%s'", table)
		}
		log.Trace().Str("table", table).Msg("data clean done")
	}
	return nil
}
