package app

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/internal/app/domain/mylibrary"
	"github.com/typical-go/typical-rest-server/internal/app/domain/mymusic"
	"github.com/typical-go/typical-rest-server/internal/app/infra"
	"github.com/typical-go/typical-rest-server/internal/app/infra/log"
	"github.com/typical-go/typical-rest-server/pkg/logruskit"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"

	// enable `/debug/vars`
	_ "expvar"

	// enable `/debug/pprof` API
	_ "net/http/pprof"
)

const (
	healthCheckPath = "/application/health"
)

type (
	app struct {
		dig.In
		Logger      *logrus.Logger
		Config      *infra.AppCfg
		Library     mylibrary.Router
		Album       mymusic.Router
		HealthCheck infra.HealthCheck
	}
)

// Start app
func Start(a app) (err error) {
	e := echo.New()
	defer shutdown(e)

	e.HideBanner = true
	e.Debug = a.Config.Debug
	e.Logger = logruskit.EchoLogger(a.Logger)

	setMiddleware(a, e)
	setRoute(a, e)

	return e.StartServer(&http.Server{
		Addr:         a.Config.Address,
		ReadTimeout:  a.Config.ReadTimeout,
		WriteTimeout: a.Config.WriteTimeout,
	})
}

func setMiddleware(a app, e *echo.Echo) {
	e.Use(middleware.Recover())
	if e.Debug {
		e.Use(log.Middleware)
	}
}

func setRoute(a app, e *echo.Echo) {
	typrest.SetRoute(e,
		&a.Library,
		&a.Album,
	)

	e.GET(healthCheckPath, a.HealthCheck.Handle)
	e.HEAD(healthCheckPath, a.HealthCheck.Handle)
	e.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))
	e.GET("/debug/*/*", echo.WrapHandler(http.DefaultServeMux))

	if a.Config.Debug {
		logrus.Debugf("Application routes:\n  %s\n\n",
			strings.Join(typrest.DumpEcho(e), "\n  "))
	}
}

func shutdown(e *echo.Echo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	e.Shutdown(ctx)
}
