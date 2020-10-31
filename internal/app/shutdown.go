package app

import (
	"context"
	"database/sql"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-go/pkg/errkit"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
	"go.uber.org/dig"
)

type (
	shutdown struct {
		dig.In
		Pg    *sql.DB `name:"pg"`
		MySQL *sql.DB `name:"mysql"`
		Cache *cachekit.Store
		Echo  *echo.Echo
	}
)

// Shutdown infra
func Shutdown(p shutdown) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errs := errkit.Errors{
		p.Pg.Close(),
		p.MySQL.Close(),
		p.Cache.Close(),
		p.Echo.Shutdown(ctx),
	}

	return errs.Unwrap()
}
