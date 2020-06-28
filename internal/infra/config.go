package infra

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-redis/redis"

	// postgres driver
	_ "github.com/lib/pq"
)

type (
	// App is application configuration
	App struct {
		Address string `envconfig:"ADDRESS" default:":8089" required:"true"`
		Debug   bool   `default:"true"`
	}
	// Redis Configuration
	Redis struct {
		Host     string `required:"true" default:"localhost"`
		Port     string `required:"true" default:"6379"`
		Password string `default:"redispass"`
		DB       int    `default:"0"`

		PoolSize           int           `envconfig:"POOL_SIZE"  default:"20" required:"true"`
		DialTimeout        time.Duration `envconfig:"DIAL_TIMEOUT" default:"5s" required:"true"`
		ReadWriteTimeout   time.Duration `envconfig:"READ_WRITE_TIMEOUT" default:"3s" required:"true"`
		IdleTimeout        time.Duration `envconfig:"IDLE_TIMEOUT" default:"5m" required:"true"`
		IdleCheckFrequency time.Duration `envconfig:"IDLE_CHECK_FREQUENCY" default:"1m" required:"true"`
		MaxConnAge         time.Duration `envconfig:"MAX_CONN_AGE" default:"30m" required:"true"`
	}
	// Pg is postgres configuration
	Pg struct {
		DBName   string `required:"true" default:"MyLibrary"`
		User     string `required:"true" default:"postgres"`
		Password string `required:"true" default:"pgpass"`
		Host     string `default:"localhost"`
		Port     string `default:"5432"`
	}
)

//
// Redis
//

func (r *Redis) connect() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%s", r.Host, r.Port),
		Password:           r.Password,
		DB:                 r.DB,
		PoolSize:           r.PoolSize,
		DialTimeout:        r.DialTimeout,
		ReadTimeout:        r.ReadWriteTimeout,
		WriteTimeout:       r.ReadWriteTimeout,
		IdleTimeout:        r.IdleTimeout,
		IdleCheckFrequency: r.IdleCheckFrequency,
		MaxConnAge:         r.MaxConnAge,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, fmt.Errorf("infra: %w", err)
	}

	return client, nil
}

//
// Pg
//

func (p *Pg) connect() (*sql.DB, error) {
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
