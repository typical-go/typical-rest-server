package typnewrelic

import (
	"fmt"

	"github.com/labstack/echo"

	nr "github.com/newrelic/go-agent"
	log "github.com/sirupsen/logrus"
)

const (
	// ErrorMessage is echo context key that contain error message to be relay as NewRelic error
	ErrorMessage = "error-message"
)

// Middleware to initiation middleware for newrelic
func Middleware(e *echo.Echo, nrApp nr.Application) {
	if nrApp != nil {
		log.Info("Use NewRelic Middleware")
		e.Use(MiddlewareFunc(nrApp))
	} else {
		log.Info("NewRelic Application not found")
	}
}

// MiddlewareFunc to capture error and panic to newrelic
func MiddlewareFunc(app nr.Application) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// start new relic transaction
			name := fmt.Sprintf("%s [%s]", ctx.Path(), ctx.Request().Method)
			txn := app.StartTransaction(name, ctx.Response().Writer, ctx.Request())
			defer txn.End()
			req := ctx.Request()
			nrContext := nr.NewContext(req.Context(), txn)
			ctx.SetRequest(req.WithContext(nrContext))
			// handle panic
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("PANIC: %v", r)
					}
					ctx.Error(err)
					txn.NoticeError(err)
					log.Error(err.Error())
				}
			}()
			// run the actual handler
			err := next(ctx)
			if err != nil {
				txn.NoticeError(err)
				log.Error(err.Error())
			}
			errMessage, ok := ctx.Get(ErrorMessage).(string)
			if ok && errMessage != "" {
				txn.NoticeError(err)
				log.Error(ErrorMessage)
			}
			return err
		}
	}
}
