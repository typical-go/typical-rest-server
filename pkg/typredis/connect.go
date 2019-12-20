package typredis

import (
	"fmt"

	"github.com/go-redis/redis"
)

func (*module) connect(cfg Config) (client *redis.Client) {
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
	return
}

func (*module) ping(client *redis.Client) (err error) {
	if err = client.Ping().Err(); err != nil {
		return fmt.Errorf("Redis: Ping: %w", err)
	}
	return
}

func (*module) disconnect(client *redis.Client) (err error) {
	if err = client.Close(); err != nil {
		return fmt.Errorf("Redis: Disconnect: %w", err)
	}
	return
}
