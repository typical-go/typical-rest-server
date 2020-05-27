package infra

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
	// Configs of infra
	Configs struct {
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

	// Params of infra
	Params struct {
		dig.In
		Pg    *sql.DB
		Redis *redis.Client
	}
)

// Connect to infra
// @ctor
func Connect(c Configs) (infras Infras, err error) {
	pg, err := typpg.Connect(c.Pg)
	if err != nil {
		return
	}

	redis, err := typredis.Connect(c.Redis)
	if err != nil {
		return
	}

	return Infras{
		Pg:    pg,
		Redis: redis,
	}, nil
}

// Disconnect from postgres server
// @dtor
func Disconnect(p Params) (err error) {
	if err = typpg.Disconnect(p.Pg); err != nil {
		return
	}
	if err = typredis.Disconnect(p.Redis); err != nil {
		return
	}
	return
}
