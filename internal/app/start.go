package app

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/internal/app/controller"
	"github.com/typical-go/typical-rest-server/internal/app/infra"
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

// Start app
func Start(
	di *dig.Container,
	cfg *infra.EchoCfg,
	e *echo.Echo,
) (err error) {
	if err := di.Invoke(SetMiddleware); err != nil {
		return err
	}
	if err := di.Invoke(SetRoute); err != nil {
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

// SetMiddleware set middleware
func SetMiddleware(e *echo.Echo) {
	e.Use(infra.LogMiddleware)
	e.Use(middleware.Recover())
}

// SetRoute set route
func SetRoute(
	e *echo.Echo,
	hc HealthCheck,
	bookCntrl controller.BookCntrl,
) {

	// set route
	echokit.SetRoute(e, &bookCntrl)

	// profiling
	e.GET(healthCheckPath, hc.Handle)
	e.HEAD(healthCheckPath, hc.Handle)
	e.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))
	e.GET("/debug/*/*", echo.WrapHandler(http.DefaultServeMux))
}
