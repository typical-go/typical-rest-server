package server

import (
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/pkg/serverkit"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"go.uber.org/dig"
)

type healthcheck struct {
	dig.In
	PG    *typpostgres.DB
	Redis *redis.Client
}

func (h *healthcheck) SetRoute(e *echo.Echo) {
	healthcheck := serverkit.NewHealthCheck()
	healthcheck.Put("postgres", h.PG.Ping)
	healthcheck.Put("redis", h.Redis.Ping().Err)

	status, message := healthcheck.Process()
	e.Any("application/health", func(ec echo.Context) (err error) {
		return ec.JSON(status, message)
	})
}
