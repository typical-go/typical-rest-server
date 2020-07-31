package app

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/typical-go/typical-rest-server/internal/infra"
	"github.com/typical-go/typical-rest-server/internal/profiler"
	"github.com/typical-go/typical-rest-server/internal/server"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"go.uber.org/dig"
)

type app struct {
	dig.In
	*infra.AppCfg
	Server   server.Router
	Profiler profiler.Router
}

func setMiddleware(e *echo.Echo) {
	e.Use(middleware.Recover())
	if e.Debug {
		e.Use(loggerMiddleware())
	}
}

func setRoute(e *echo.Echo, a *app) error {
	return echokit.SetRoute(e,
		&a.Server,
		&a.Profiler,
	)
}

func shutdown(e *echo.Echo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	e.Shutdown(ctx)
}
