package profiler

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/pkg/serverkit"
	"go.uber.org/dig"
)

type Controller struct {
	dig.In
	PG    *sql.DB
	Redis *redis.Client
}

func (h *Controller) SetRoute(e *echo.Echo) {
	e.Any("application/health", h.healthCheck)
}

func (h *Controller) healthCheck(ec echo.Context) (err error) {
	healthcheck := serverkit.NewHealthCheck()
	healthcheck.Put("postgres", h.PG.Ping)
	healthcheck.Put("redis", h.Redis.Ping().Err)

	status, message := healthcheck.Process()
	return ec.JSON(status, message)
}
