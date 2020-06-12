package infra

import (
	"database/sql"
	"fmt"

	"github.com/go-redis/redis"
	"go.uber.org/dig"

	// postgres driver
	_ "github.com/lib/pq"
)

type (
	// Configs of infra
	Configs struct {
		dig.In
		Pg    *Pg
		Redis *Redis
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
	pg, err := connectPg(c.Pg)
	if err != nil {
		return
	}

	redis, err := connectRedis(c.Redis)
	if err != nil {
		return
	}

	return Infras{
		Pg:    pg,
		Redis: redis,
	}, nil
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

func connectRedis(cfg *Redis) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:           cfg.Password,
		DB:                 cfg.DB,
		PoolSize:           cfg.PoolSize,
		DialTimeout:        cfg.DialTimeout,
		ReadTimeout:        cfg.ReadWriteTimeout,
		WriteTimeout:       cfg.ReadWriteTimeout,
		IdleTimeout:        cfg.IdleTimeout,
		IdleCheckFrequency: cfg.IdleCheckFrequency,
		MaxConnAge:         cfg.MaxConnAge,
	})

	if err = client.Ping().Err(); err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}

	return client, nil
}

func connectPg(cfg *Pg) (*sql.DB, error) {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		err = fmt.Errorf("postgres: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	return db, nil
}
