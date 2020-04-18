package server

import (
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/pkg/serverkit"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"go.uber.org/dig"
)

type profiler struct {
	dig.In
	PG    *typpostgres.DB
	Redis *redis.Client
}

func (h *profiler) SetRoute(e *echo.Echo) {
	e.Any("application/health", h.healthCheck)
}

func (h *profiler) healthCheck(ec echo.Context) (err error) {
	healthcheck := serverkit.NewHealthCheck()
	healthcheck.Put("postgres", h.PG.Ping)
	healthcheck.Put("redis", h.Redis.Ping().Err)

	status, message := healthcheck.Process()
	return ec.JSON(status, message)
}
