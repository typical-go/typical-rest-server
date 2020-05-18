package server

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/typpg"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
	"go.uber.org/dig"

	// postgres driver
	_ "github.com/lib/pq"
)

type (
	// InfraConfigs is config collection
	InfraConfigs struct {
		dig.In
		Pg    *typpg.Config
		Redis *typredis.Config
	}

	// Infras is list of infra to be provide in dependency-injection
	Infras struct {
		dig.Out
		Pg    *sql.DB
		Redis *redis.Client
	}

	// InfraParams is infra as parameter
	InfraParams struct {
		dig.In
		Pg    *sql.DB
		Redis *redis.Client
	}
)

// Connect to infra
// @ctor
func Connect(c InfraConfigs) (infras Infras, err error) {
	var (
		pg    *sql.DB
		redis *redis.Client
	)

	if pg, err = typpg.Connect(c.Pg); err != nil {
		return
	}

	if redis, err = typredis.Connect(c.Redis); err != nil {
		return
	}

	return Infras{
		Pg:    pg,
		Redis: redis,
	}, nil
}

// Disconnect from postgres server
// @dtor
func Disconnect(p InfraParams) (err error) {
	if err = typpg.Disconnect(p.Pg); err != nil {
		return
	}
	if err = typredis.Disconnect(p.Redis); err != nil {
		return
	}
	return
}
