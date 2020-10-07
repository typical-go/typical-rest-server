package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/internal/app/domain/mylibrary"
	"github.com/typical-go/typical-rest-server/internal/app/domain/mymusic"
	"github.com/typical-go/typical-rest-server/internal/app/infra"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"

	// enable `/debug/vars`
	_ "expvar"

	// enable `/debug/pprof` API
	_ "net/http/pprof"
)

type (
	app struct {
		dig.In
		Config      *infra.AppCfg
		Library     mylibrary.Router
		Album       mymusic.Router
		HealthCheck infra.HealthCheck
	}
)

const (
	healthCheckPath = "/application/health"
)

// Start app
func Start(a app) (err error) {
	e := echo.New()
	defer Shutdown(e)

	e.HideBanner = true
	e.Debug = a.Config.Debug

	a.SetLoggger(e)
	a.SetMiddleware(e)
	a.SetRoute(e)

	return e.StartServer(&http.Server{
		Addr:         a.Config.Address,
		ReadTimeout:  a.Config.ReadTimeout,
		WriteTimeout: a.Config.WriteTimeout,
	})
}

// Shutdown app
func Shutdown(e *echo.Echo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	e.Shutdown(ctx)
}

//
// app
//

// SetMiddleware set middleware to the app
func (a app) SetMiddleware(e *echo.Echo) {
	e.Use(middleware.Recover())
	if e.Debug {
		e.Use(loggerMiddleware())
	}
}

// SetRoute set route the app
func (a app) SetRoute(e *echo.Echo) {
	typrest.SetRoute(e,
		&a.Library,
		&a.Album,
	)

	e.GET(healthCheckPath, a.HealthCheck.Handle)
	e.HEAD(healthCheckPath, a.HealthCheck.Handle)
	e.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))
	e.GET("/debug/*/*", echo.WrapHandler(http.DefaultServeMux))

	if a.Config.Debug {
		var routePaths []string
		for _, route := range e.Routes() {
			routePaths = append(routePaths, fmt.Sprintf("  %s\t%s", route.Path, route.Method))
		}
		sort.Strings(routePaths)
		logrus.Debugf("Application routes:\n%s\n\n", strings.Join(routePaths, "\n"))
	}
}

// SetLogger set logger to the app
func (a app) SetLoggger(e *echo.Echo) {
	logger := logrus.StandardLogger()     // NOTE: use standard logger for global use
	e.Logger = typrest.WrapLogrus(logger) // NOTE: setup echo logger
	log.SetOutput(logger.Writer())        // NOTE: std golang log use same output writer with logrus

	if a.Config.Debug {
		logger.SetLevel(logrus.DebugLevel)
		logger.SetFormatter(&logrus.TextFormatter{})
	} else {
		logger.SetLevel(logrus.WarnLevel)
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
}

// loggerMiddleware log every request
func loggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			bytesIn := req.Header.Get(echo.HeaderContentLength)

			logrus.WithFields(map[string]interface{}{
				"status":    res.Status,
				"latency":   stop.Sub(start).String(),
				"bytes_in":  bytesIn,
				"bytes_out": strconv.FormatInt(res.Size, 10),
			}).Info(fmt.Sprintf("%s %s", req.Method, req.RequestURI))
			return nil
		}
	}
}
