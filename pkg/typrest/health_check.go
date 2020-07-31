package typrest

import (
	"net/http"
)

// HealthStatusOK is health status ok
var HealthStatusOK = "OK"

type (
	// HealthCheck to handle health-check
	HealthCheck map[string]func() error
)

// Result of HealthCheck
func (h HealthCheck) Result() (status int, message map[string]string) {
	status = http.StatusOK
	message = make(map[string]string)

	for name, fn := range h {
		if err := fn(); err != nil {
			message[name] = err.Error()
			status = http.StatusServiceUnavailable
		} else {
			message[name] = HealthStatusOK
		}
	}
	return
}
