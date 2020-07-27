package app

import (
	"context"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/internal/infra"
	"github.com/typical-go/typical-rest-server/internal/profiler"
	"github.com/typical-go/typical-rest-server/internal/server"

	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"github.com/typical-go/typical-rest-server/pkg/echologrus"

	"go.uber.org/dig"
)

type app struct {
	dig.In
	*infra.AppCfg
	Server   server.Router
	Profiler profiler.Router
}

var _ echokit.Router = (*app)(nil)

// Main function to run server
func Main(a app) (err error) {
	e := echo.New()
	defer shutdown(e)

	e.HideBanner = true
	a.initLogger(e)
	a.Middleware(e)
	a.Route(e)

	return e.Start(a.Address)
}

func (a app) Middleware(e echokit.Server) {
	e.Use(middleware.Recover())
}

func (a app) Route(e echokit.Server) (err error) {
	return echokit.SetRoute(e,
		&a.Server,
		&a.Profiler,
	)
}

func (a app) initLogger(e *echo.Echo) {
	e.Logger = echologrus.Logger{Logger: log.StandardLogger()}
	e.Debug = a.Debug
	if e.Debug {
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{})
		e.Use(echologrus.HookWithConfig(echologrus.Config{
			IncludeRequestBodies:  true,
			IncludeResponseBodies: true,
		}))
	} else {
		log.SetLevel(log.WarnLevel)
		log.SetFormatter(&log.JSONFormatter{})
		e.Use(echologrus.HookWithConfig(echologrus.Config{}))
	}
}

func shutdown(e *echo.Echo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	e.Shutdown(ctx)
}
