package restserver

import (
	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
	"github.com/typical-go/typical-rest-server/restserver/config"
	"github.com/typical-go/typical-rest-server/restserver/controller"

	"go.uber.org/dig"
)

type app struct {
	dig.In
	*typserver.Server
	config.Config
	controller.BookCntrl

	Postgres *typpostgres.DB
	Redis    *redis.Client
}

func startServer(a app) (err error) {
	a.SetDebug(a.Debug)

	// health check
	a.PutHealthChecker("postgres", a.Postgres.Ping)
	a.PutHealthChecker("redis", a.Redis.Ping().Err)

	// set middleware
	// a.Use(middleware.Recover()) // TODO: uncomment when

	// register controller
	a.Register(&a.BookCntrl)

	return a.Start(a.Address)
}
