package dbkit

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

// SetCache to set cache
func SetCache(client *redis.Client, key string, val interface{}, exp time.Duration) error {
	data, _ := json.Marshal(val)
	return client.Set(key, data, exp).Err()
}

// GetCache to get cache
func GetCache(client *redis.Client, key string, val interface{}) (err error) {
	data, err := client.Get(key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, val)
}
