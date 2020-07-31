package app

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func setLogLevel(logger *logrus.Logger, debug bool) {
	if debug {
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
