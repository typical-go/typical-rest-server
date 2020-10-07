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
		Pg    *sql.DB `name:"pg"`
		MySQL *sql.DB `name:"mysql"`
		Redis *redis.Client
	}
	setupParam struct {
		dig.In
		PgCfg    *PostgresCfg
		MysqlCfg *MySQLCfg
		RedisCfg *RedisCfg
	}
	teardownParam struct {
		dig.In
		Pg    *sql.DB `name:"pg"`
		MySQL *sql.DB `name:"mysql"`
		Redis *redis.Client
	}
)

// Setup infra
// @ctor
func Setup(p setupParam) Infra {
	return Infra{
		Pg:    p.PgCfg.createConn(),
		MySQL: p.MysqlCfg.createConn(),
		Redis: p.RedisCfg.createClient(),
	}
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
