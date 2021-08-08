package app

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/dig"

	"time"

	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-go/pkg/errkit"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
)

// Shutdown infra
func Shutdown(p struct {
	dig.In
	Pg    *sql.DB `name:"pg"`
	Cache *cachekit.Store
	Echo  *echo.Echo
}) error {

	fmt.Printf("Shutdown at %s", time.Now().String())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errs := errkit.Errors{
		p.Pg.Close(),
		p.Cache.Close(),
		p.Echo.Shutdown(ctx),
	}

	return errs.Unwrap()
}
