package dockerrx

import (
	"errors"
	"os"

	"github.com/typical-go/typical-rest-server/pkg/typdocker"
)

type (
	// RedisWithEnv same with redis with env parameter
	RedisWithEnv struct {
		Version     string
		Name        string
		Image       string
		PasswordEnv string
		PortEnv     string
	}
)

var _ (typdocker.Composer) = (*RedisWithEnv)(nil)

// ComposeV3 for docker-compose
func (r *RedisWithEnv) ComposeV3() (*typdocker.Recipe, error) {
	if r.PasswordEnv == "" {
		return nil, errors.New("redis-with-env: missing PasswordEnv")
	}
	if r.PortEnv == "" {
		return nil, errors.New("redis-with-env: redis: missing PortEnv")
	}
	redis := Redis{
		Version:  r.Version,
		Name:     r.Name,
		Image:    r.Image,
		Password: os.Getenv(r.PasswordEnv),
		Port:     os.Getenv(r.PortEnv),
	}
	return redis.ComposeV3()
}
