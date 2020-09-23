package profiler

import (
	"database/sql"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"
)

type (
	// HealthCheck for profiler
	HealthCheck struct {
		dig.In
		PG    *sql.DB
		Redis *redis.Client
	}
	// HealthCheckResponse healthcheck response
	HealthCheckResponse struct {
		AppName    string            `json:"name"`
		AppVersion string            `json:"version"`
		Status     map[string]string `json:"status"`
	}
)

var _ typrest.Router = (*HealthCheck)(nil)

// SetRoute to profiler api
func (h *HealthCheck) SetRoute(e typrest.Server) {
	e.Any("application/health", h.handle)
}

func (h *HealthCheck) handle(ec echo.Context) error {
	healthy, detail := typrest.HealthStatus(typrest.HealthMap{
		"postgres": h.PG.Ping,
		"redis":    h.Redis.Ping().Err,
	})

	status := http.StatusOK
	if !healthy {
		status = http.StatusServiceUnavailable
	}

	return ec.JSON(status, HealthCheckResponse{
		AppName:    typgo.AppName,
		AppVersion: typgo.AppVersion,
		Status:     detail,
	})
}
