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
)

// Connect to infra
// @ctor
func Connect(pgCfg *PostgresCfg, redisCfg *RedisCfg) (infras Infra, err error) {
	pg, err := pgCfg.connect()
	if err != nil {
		return
	}
	redis, err := redisCfg.connect()
	if err != nil {
		return
	}
	return Infra{Pg: pg, Redis: redis}, nil
}

// Disconnect from postgres server
// @dtor
func Disconnect(pg *sql.DB, redis *redis.Client) error {
	if err := pg.Close(); err != nil {
		return err
	}
	if err := redis.Close(); err != nil {
		return err
	}
	return nil
}
