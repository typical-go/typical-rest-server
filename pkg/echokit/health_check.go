package echokit

import (
	"net/http"

	"github.com/labstack/echo"
)

// OK is status when no error
const OK = "OK"

// HealthCheck is key-value store that contain health check information
type HealthCheck map[string]string

// NewHealthCheck return new instance of HealthCheck
func NewHealthCheck() HealthCheck {
	return HealthCheck(make(map[string]string))
}

// Add name and error to register as heath check
func (c HealthCheck) Add(name string, err error) HealthCheck {
	status := OK
	if err != nil {
		status = err.Error()
	}

	c[name] = status
	return c
}

// NotOK return true is some error registered
func (c HealthCheck) NotOK() bool {
	for _, value := range c {
		if value != OK {
			return true
		}
	}
	return false
}

// Send healthcheck response
func (c HealthCheck) Send(ctx echo.Context) error {
	var status int
	if c.NotOK() {
		status = http.StatusServiceUnavailable
	} else {
		status = http.StatusOK
	}
	return ctx.JSON(status, c)
}
