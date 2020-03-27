package server

import (
	"github.com/go-redis/redis"
	"github.com/labstack/echo/middleware"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
	"github.com/typical-go/typical-rest-server/server/config"
	"github.com/typical-go/typical-rest-server/server/controller"

	"go.uber.org/dig"
)

var (
	configName = "APP"
)

type server struct {
	dig.In
	*typserver.Server
	*config.Config

	controller.BookCntrl

	Postgres *typpostgres.DB
	Redis    *redis.Client
}

// Main function to run server
func Main(s server) (err error) {
	s.SetLogger(s.Debug)

	// health check
	s.PutHealthChecker("postgres", s.Postgres.Ping)
	s.PutHealthChecker("redis", s.Redis.Ping().Err)

	// set middleware
	s.Use(middleware.Recover())

	// register controller
	s.Register(&s.BookCntrl)

	return s.Start(s.Address)
}

// Configuration of server
func Configuration() *typcfg.Configuration {
	return typcfg.NewConfiguration(configName, &config.Config{
		Debug: true,
	})
}
