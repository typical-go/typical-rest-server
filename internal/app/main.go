package app

import (
	"context"
	"time"

	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/typical-go/typical-rest-server/internal/app/infra"
	"github.com/typical-go/typical-rest-server/internal/app/profiler"
	"github.com/typical-go/typical-rest-server/internal/server"
	"github.com/typical-go/typical-rest-server/pkg/echokit"

	"go.uber.org/dig"
)

type app struct {
	dig.In
	*infra.App
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
	e.Logger = logrusmiddleware.Logger{Logger: log.StandardLogger()}
	e.Debug = a.Debug
	if e.Debug {
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{})
		e.Use(logrusmiddleware.HookWithConfig(logrusmiddleware.Config{
			IncludeRequestBodies:  true,
			IncludeResponseBodies: true,
		}))
	} else {
		log.SetLevel(log.WarnLevel)
		log.SetFormatter(&log.JSONFormatter{})
		e.Use(logrusmiddleware.HookWithConfig(logrusmiddleware.Config{}))
	}
}

func shutdown(e *echo.Echo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	e.Shutdown(ctx)
}
