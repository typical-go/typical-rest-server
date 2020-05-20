package app

import (
	"context"
	"time"

	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/typical-go/typical-rest-server/internal/config"
	"github.com/typical-go/typical-rest-server/internal/profiler"
	"github.com/typical-go/typical-rest-server/internal/server/controller"

	"go.uber.org/dig"
)

type app struct {
	dig.In
	Config    *config.Config
	BookCntrl controller.BookCntrl
	Profiler  profiler.Controller
}

// Main function to run server
func Main(a app) (err error) {
	e := echo.New()
	defer shutdown(e)

	e.HideBanner = true
	initLogger(e, a.Config.Debug)

	// set middleware
	e.Use(middleware.Recover())

	a.BookCntrl.SetRoute(e)
	a.Profiler.SetRoute(e)

	return e.Start(a.Config.Address)
}

func initLogger(e *echo.Echo, debug bool) {
	e.Logger = logrusmiddleware.Logger{Logger: log.StandardLogger()}
	e.Debug = debug
	if debug {
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
