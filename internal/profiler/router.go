package profiler

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"
)

// Router for profiler
type Router struct {
	dig.In
	PG    *sql.DB
	Redis *redis.Client
}

var _ typrest.Router = (*Router)(nil)

// SetRoute to profiler api
func (h *Router) SetRoute(e typrest.Server) error {
	e.Any("application/health", h.healthCheck)
	return nil
}

func (h *Router) healthCheck(ec echo.Context) error {
	hc := typrest.HealthCheck{
		"postgres": h.PG.Ping,
		"redis":    h.Redis.Ping().Err,
	}
	status, message := hc.Result()
	return ec.JSON(status, message)
}
