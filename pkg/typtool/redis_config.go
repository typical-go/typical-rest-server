package typtool

import "os"

var DefaultRedisEnvKeys = &RedisEnvKeys{
	Host: "HOST",
	Port: "PORT",
	Pass: "PASS",
}

type (
	RedisConfig struct {
		Host string
		Port string
		Pass string
	}
	RedisEnvKeys RedisConfig
)

func RedisEnvKeysWithPrefix(prefix string) *RedisEnvKeys {
	return &RedisEnvKeys{
		Host: prefix + "_" + DefaultRedisEnvKeys.Host,
		Port: prefix + "_" + DefaultRedisEnvKeys.Port,
		Pass: prefix + "_" + DefaultRedisEnvKeys.Pass,
	}
}

func (r *RedisEnvKeys) Config() *RedisConfig {
	return &RedisConfig{
		Host: os.Getenv(r.Host),
		Port: os.Getenv(r.Port),
		Pass: os.Getenv(r.Pass),
	}
}
