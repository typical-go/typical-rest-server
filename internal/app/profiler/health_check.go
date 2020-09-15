package profiler

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
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
)

var _ typrest.Router = (*HealthCheck)(nil)

// SetRoute to profiler api
func (h *HealthCheck) SetRoute(e typrest.Server) {
	e.Any("application/health", h.handle)
}

func (h *HealthCheck) handle(ec echo.Context) error {
	hc := typrest.HealthCheck{
		"postgres": h.PG.Ping,
		"redis":    h.Redis.Ping().Err,
	}
	status, message := hc.Result()
	return ec.JSON(status, message)
}
