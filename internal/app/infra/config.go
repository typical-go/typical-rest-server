package infra

import (
	"database/sql"
	"fmt"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"

	// postgres driver
	_ "github.com/lib/pq"
)

type (
	// AppCfg application configuration
	// @envconfig (prefix:"APP")
	AppCfg struct {
		Address string `envconfig:"ADDRESS" default:":8089" required:"true"`
		Debug   bool   `envconfig:"DEBUG" default:"true"`
	}
	// RedisCfg redis onfiguration
	// @envconfig (prefix:"REDIS")
	RedisCfg struct {
		Host     string `envconfig:"HOST" required:"true" default:"localhost"`
		Port     string `envconfig:"PORT" required:"true" default:"6379"`
		Password string `envconfig:"PASSWORD" default:"redispass"`
	}
	// PostgresCfg postgres configuration
	// @envconfig (prefix:"PG")
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

func (r *RedisCfg) createClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Password: r.Password,
	})

	if err := client.Ping().Err(); err != nil {
		log.Fatalf("infra: %s", err.Error())
	}

	return client
}

//
// PostgresCfg
//

func (p *PostgresCfg) createConn() *sql.DB {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.DBUser, p.DBPass, p.Host, p.Port, p.DBName,
	)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatalf("infra: %s", err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("infra: %s", err.Error())
	}

	return db
}
