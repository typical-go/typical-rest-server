package infra

import (
	"database/sql"

	"github.com/go-redis/redis"
	"go.uber.org/dig"
)

type (
	// Configs of infra
	Configs struct {
		dig.In
		Pg    *PostgresCfg
		Redis *RedisCfg
	}
	// Infras is list of infra to be provide in dependency-injection
	Infras struct {
		dig.Out
		Pg    *sql.DB
		Redis *redis.Client
	}
	// Params of infra
	Params struct {
		dig.In
		Pg    *sql.DB
		Redis *redis.Client
	}
)

// Connect to infra
// @ctor
func Connect(c Configs) (infras Infras, err error) {
	pg, err := c.Pg.connect()
	if err != nil {
		return
	}
	redis, err := c.Redis.connect()
	if err != nil {
		return
	}
	return Infras{Pg: pg, Redis: redis}, nil
}

// Disconnect from postgres server
// @dtor
func Disconnect(p Params) error {
	if err := p.Pg.Close(); err != nil {
		return err
	}
	if err := p.Redis.Close(); err != nil {
		return err
	}
	return nil
}
