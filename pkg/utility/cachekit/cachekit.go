package cachekit

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

// Set the cache
func Set(client *redis.Client, key string, val interface{}, exp time.Duration) error {
	data, _ := json.Marshal(val)
	return client.Set(key, data, exp).Err()
}

// Get the cache
func Get(client *redis.Client, key string, val interface{}) (err error) {
	data, err := client.Get(key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, val)
}
