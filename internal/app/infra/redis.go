package infra

import (
	"fmt"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

func createRedisClient(r *RedisCfg) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Password: r.Password,
	})

	if err := client.Ping().Err(); err != nil {
		log.Fatalf("redis: %s", err.Error())
	}

	return client
}
