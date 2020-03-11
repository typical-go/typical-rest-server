package restserver

import (
	"github.com/go-redis/redis"
	"github.com/labstack/echo/middleware"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
	"github.com/typical-go/typical-rest-server/restserver/config"
	"github.com/typical-go/typical-rest-server/restserver/controller"

	"go.uber.org/dig"
)

type server struct {
	dig.In
	*typserver.Server
	config.Config
	controller.BookCntrl
	Postgres *typpostgres.DB
	Redis    *redis.Client
}

func startServer(s server) (err error) {
	s.SetLogger(s.Debug)

	// health check
	s.PutHealthChecker("postgres", s.Postgres.Ping)
	s.PutHealthChecker("redis", s.Redis.Ping().Err)

	// set middleware
	s.Use(middleware.Recover()) // TODO: uncomment when

	// register controller
	s.Register(&s.BookCntrl)

	return s.Start(s.Address)
}
