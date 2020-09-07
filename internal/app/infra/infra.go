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
	setupParam struct {
		dig.In
		PgCfg    *PostgresCfg
		RedisCfg *RedisCfg
	}
	teardownParam struct {
		dig.In
		Pg    *sql.DB
		Redis *redis.Client
	}
)

// Setup infra
// @ctor
func Setup(p setupParam) (infras Infra, err error) {
	pg, err := p.PgCfg.createConn()
	if err != nil {
		return
	}
	redis, err := p.RedisCfg.createClient()
	if err != nil {
		return
	}
	return Infra{Pg: pg, Redis: redis}, nil
}

// Teardown infra
// @dtor
func Teardown(p teardownParam) error {
	if err := p.Pg.Close(); err != nil {
		return err
	}
	if err := p.Redis.Close(); err != nil {
		return err
	}
	return nil
}
