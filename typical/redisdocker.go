package typical

import (
	"os"

	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/pkg/dockerrx"
)

type redisDocker struct {
	name string
}

var _ (typdocker.Composer) = (*redisDocker)(nil)

func (r *redisDocker) Compose() (*typdocker.Recipe, error) {
	redis := dockerrx.Redis{
		Version:  typdocker.V3,
		Name:     r.name,
		Image:    "redis:4.0.5-alpine",
		Password: os.Getenv("REDIS_PASSWORD"),
		Port:     os.Getenv("REDIS_PORT"),
	}
	return redis.Compose()
}
