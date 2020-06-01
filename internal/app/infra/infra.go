package infra

import (
	"database/sql"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/typpg"
	"go.uber.org/dig"

	// postgres driver
	_ "github.com/lib/pq"
)

type (
	// Configs of infra
	Configs struct {
		dig.In
		Pg    *typpg.Config
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
	pg, err := typpg.Connect(c.Pg)
	if err != nil {
		return
	}

	redis, err := openRedis(c.Redis)
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
func Disconnect(p Params) (err error) {
	if err = typpg.Disconnect(p.Pg); err != nil {
		return
	}
	if err = p.Redis.Close(); err != nil {
		return
	}
	return
}

func openRedis(cfg *Redis) (client *redis.Client, err error) {
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
