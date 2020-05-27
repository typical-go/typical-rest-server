package profiler

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"go.uber.org/dig"
)

type Controller struct {
	dig.In
	PG    *sql.DB
	Redis *redis.Client
}

func (h *Controller) SetRoute(e *echo.Echo) {
	hc := echokit.HealthCheck{
		"postgres": h.PG.Ping,
		"redis":    h.Redis.Ping().Err,
	}

	e.Any("application/health", hc.JSON)
}
