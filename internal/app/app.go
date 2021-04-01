package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/typical-go/typical-go/pkg/errkit"
	"github.com/typical-go/typical-rest-server/internal/app/infra"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"go.uber.org/dig"

	// enable `/debug/vars`
	_ "expvar"

	// enable `/debug/pprof` API
	_ "net/http/pprof"
)

// Start app
func Start(
	di *dig.Container,
	cfg *infra.AppCfg,
	e *echo.Echo,
) (err error) {
	if err := di.Invoke(SetServer); err != nil {
		return err
	}
	if err := di.Invoke(SetProfiler); err != nil {
		return err
	}
	if cfg.Debug {
		routes := echokit.DumpEcho(e)
		logrus.Debugf("Print routes:\n  %s\n\n", strings.Join(routes, "\n  "))
	}
	return e.StartServer(&http.Server{
		Addr:         cfg.Address,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
}

// Shutdown infra
func Shutdown(p struct {
	dig.In
	Pg    *sql.DB `name:"pg"`
	MySQL *sql.DB `name:"mysql"`
	Cache *cachekit.Store
	Echo  *echo.Echo
}) error {

	fmt.Printf("Shutdown at %s", time.Now().String())
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
