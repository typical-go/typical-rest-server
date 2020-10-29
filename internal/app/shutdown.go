package app

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
	"go.uber.org/multierr"
)

type (
	shutdown struct {
		dig.In
		Pg    *sql.DB `name:"pg"`
		MySQL *sql.DB `name:"mysql"`
		Redis *redis.Client
		Echo  *echo.Echo
	}
)

// Shutdown infra
func Shutdown(p shutdown) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return multierr.Combine(
		p.Pg.Close(),
		p.MySQL.Close(),
		p.Redis.Close(),
		p.Echo.Shutdown(ctx),
	)
}
