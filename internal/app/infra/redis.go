package infra

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

// NewRedisClient return new instance of redis client
// @ctor
func NewRedisClient(r *RedisCfg) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Password: r.Password,
	})
	if err := client.Ping().Err(); err != nil {
		logrus.Fatalf("redis: %s", err.Error())
	}
	return client
}
