package app

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
)

func initLogger(e *echo.Echo) {
	echoLog := echokit.WrapLogrus(logrus.StandardLogger())
	e.Logger = echoLog

	if e.Debug {
		e.Use(Logger(echoLog))
		logrus.SetLevel(log.DebugLevel)
		logrus.SetFormatter(&log.TextFormatter{})
	} else {
		logrus.SetLevel(log.WarnLevel)
		logrus.SetFormatter(&log.JSONFormatter{})
	}
}

// Logger to log every request
func Logger(l *echokit.EchoLogrus) echo.MiddlewareFunc {
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

			l.Logger.WithFields(map[string]interface{}{
				"status":    res.Status,
				"latency":   stop.Sub(start).String(),
				"bytes_in":  bytesIn,
				"bytes_out": strconv.FormatInt(res.Size, 10),
			}).Info(fmt.Sprintf("%s %s", req.Method, req.RequestURI))
			return nil
		}
	}
}
