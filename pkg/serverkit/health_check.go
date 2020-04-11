package serverkit

import (
	"net/http"
)

// HealthCheck to handle health-check
type HealthCheck struct {
	m map[string]PingFn
}

// PingFn is ping function
type PingFn func() error

// NewHealthCheck return new instance of HealthCheck
func NewHealthCheck() *HealthCheck {
	return &HealthCheck{
		m: make(map[string]PingFn),
	}
}

// Put name and ping function
func (h *HealthCheck) Put(name string, fn PingFn) {
	h.m[name] = fn
}

// Process the ping function
func (h *HealthCheck) Process() (status int, message map[string]string) {
	status = http.StatusOK
	message = make(map[string]string)

	for name, fn := range h.m {
		if err := fn(); err != nil {
			message[name] = err.Error()
			status = http.StatusServiceUnavailable
		} else {
			message[name] = "OK"
		}
	}

	return
}
