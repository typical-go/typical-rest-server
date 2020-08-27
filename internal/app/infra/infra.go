package infra

import (
	"database/sql"

	"github.com/go-redis/redis"
	"go.uber.org/dig"
)

type (
	// Infra infrastructure for the project
	Infra struct {
		dig.Out
		Pg    *sql.DB
		Redis *redis.Client
	}
	connect struct {
		dig.In
		PgCfg    *PostgresCfg
		RedisCfg *RedisCfg
	}
	disconnect struct {
		dig.In
		Pg    *sql.DB
		Redis *redis.Client
	}
)

// Connect to infra
// @ctor
func Connect(c connect) (infras Infra, err error) {
	pg, err := c.PgCfg.createConn()
	if err != nil {
		return
	}
	redis, err := c.RedisCfg.createClient()
	if err != nil {
		return
	}
	return Infra{Pg: pg, Redis: redis}, nil
}

// Disconnect from postgres server
// @dtor
func Disconnect(d disconnect) error {
	if err := d.Pg.Close(); err != nil {
		return err
	}
	if err := d.Redis.Close(); err != nil {
		return err
	}
	return nil
}
