package typredis

import (
	"fmt"

	"github.com/go-redis/redis"
)

// Connect redis server
func Connect(cfg *Config) (client *redis.Client, err error) {
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

// Disconnect from redis server
func Disconnect(client *redis.Client) (err error) {
	if err = client.Close(); err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	return
}
