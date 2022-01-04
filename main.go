package main

// This file is here just to make `go get` easier

import (
	_ "github.com/jmoiron/sqlx"
	_ "github.com/labstack/echo/v4"
	_ "github.com/rs/zerolog"
	_ "github.com/urfave/cli/v2"
)
