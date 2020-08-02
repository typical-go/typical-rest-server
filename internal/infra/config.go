package infra

import (
	"database/sql"
	"fmt"

	"github.com/go-redis/redis"

	// postgres driver
	_ "github.com/lib/pq"
)

type (
	// AppCfg is application configuration
	// @app-cfg (prefix:"APP")
	AppCfg struct {
		Address string `envconfig:"ADDRESS" default:":8089" required:"true"`
		Debug   bool   `default:"true"`
	}
	// RedisCfg Configuration
	// @app-cfg (prefix:"REDIS")
	RedisCfg struct {
		Host     string `required:"true" default:"localhost"`
		Port     string `required:"true" default:"6379"`
		Password string `default:"redispass"`
	}
	// PostgresCfg is postgres configuration
	// @app-cfg (prefix:"PG")
	PostgresCfg struct {
		DBName   string `required:"true" default:"MyLibrary"`
		User     string `required:"true" default:"postgres"`
		Password string `required:"true" default:"pgpass"`
		Host     string `default:"localhost"`
		Port     string `default:"5432"`
	}
)

//
// RedisCfg
//

func (r *RedisCfg) connect() (*redis.Client, error) {
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

func (p *PostgresCfg) connect() (*sql.DB, error) {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.DBName)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("infra: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("infra: %w", err)
	}
	return db, nil
}
