package infra

import (
	"database/sql"
	"fmt"

	"github.com/go-redis/redis"

	// postgres driver
	_ "github.com/lib/pq"
)

type (
	// AppCfg application configuration
	// @app-cfg (prefix:"APP")
	AppCfg struct {
		Address string `envconfig:"ADDRESS" default:":8089" required:"true"`
		Debug   bool   `envconfig:"DEBUG" default:"true"`
	}
	// RedisCfg redis onfiguration
	// @app-cfg (prefix:"REDIS")
	RedisCfg struct {
		Host     string `envconfig:"HOST" required:"true" default:"localhost"`
		Port     string `envconfig:"PORT" required:"true" default:"6379"`
		Password string `envconfig:"PASSWORD" default:"redispass"`
	}
	// PostgresCfg postgres configuration
	// @app-cfg (prefix:"PG")
	PostgresCfg struct {
		DBName string `envconfig:"DBNAME" required:"true" default:"MyLibrary"`
		DBUser string `envconfig:"DBUSER" required:"true" default:"pguser"`
		DBPass string `envconfig:"DBPASS" default:"pgpass"`
		Host   string `envconfig:"HOST" required:"true" default:"localhost"`
		Port   string `envconfig:"PORT" required:"true" default:"5432"`
	}
)

//
// RedisCfg
//

func (r *RedisCfg) createClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Password: r.Password,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, fmt.Errorf("infra: %w", err)
	}

	return client, nil
}

//
// PostgresCfg
//

func (p *PostgresCfg) createConn() (*sql.DB, error) {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.DBUser, p.DBPass, p.Host, p.Port, p.DBName,
	)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("infra: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("infra: %w", err)
	}
	return db, nil
}
