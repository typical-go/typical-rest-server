package server

import (
	"context"
	"time"

	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/internal/config"
	"github.com/typical-go/typical-rest-server/internal/profiler"
	"github.com/typical-go/typical-rest-server/internal/server/controller"

	"go.uber.org/dig"
)

var (
	configName = "APP"
)

type server struct {
	dig.In
	*config.Config

	BookCntrl controller.BookCntrl
	Profiler  profiler.Controller
}

// Main function to run server
func Main(s server) (err error) {
	e := echo.New()
	defer shutdown(e)

	e.HideBanner = true
	initLogger(e, s.Debug)

	// set middleware
	e.Use(middleware.Recover())

	s.BookCntrl.SetRoute(e)
	s.Profiler.SetRoute(e)

	return e.Start(s.Address)
}

// Configuration of server
func Configuration() *typgo.Configuration {
	return &typgo.Configuration{
		Name: configName,
		Spec: &config.Config{
			Debug: true,
		},
	}
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
