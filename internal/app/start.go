package app

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/internal/app/infra"
	"github.com/typical-go/typical-rest-server/internal/app/infra/log"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
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
		*echo.Echo
		Config      *infra.AppCfg
		Router      Router
		HealthCheck HealthCheck
	}
)

// Start app
func Start(a app) (err error) {
	setMiddleware(a)
	setRoute(a)
	setProfiler(a)

	if a.Config.Debug {
		routes := echokit.DumpEcho(a.Echo)
		logrus.Debugf("Print routes:\n  %s\n\n", strings.Join(routes, "\n  "))
	}

	return a.StartServer(&http.Server{
		Addr:         a.Config.Address,
		ReadTimeout:  a.Config.ReadTimeout,
		WriteTimeout: a.Config.WriteTimeout,
	})
}

func setRoute(a app) {
	echokit.SetRoute(a, &a.Router)
}

func setMiddleware(a app) {
	a.Use(log.Middleware)
	a.Use(middleware.Recover())
}

func setProfiler(a app) {
	a.GET(healthCheckPath, a.HealthCheck.Handle)
	a.HEAD(healthCheckPath, a.HealthCheck.Handle)
	a.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))
	a.GET("/debug/*/*", echo.WrapHandler(http.DefaultServeMux))
}
