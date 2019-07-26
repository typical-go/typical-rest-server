package typical

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/dbtool"
	"github.com/typical-go/typical-rest-server/config"
)

// TODO: to be generated

type Config struct {
	App *config.AppConfig
	Pg  *config.PostgresConfig
}

func init() {
	Context.AddConstructor(func() (cfg Config, err error) {
		err = envconfig.Process("", &cfg)
		return
	})
	Context.AddConstructor(func(cfg Config) *config.AppConfig {
		return cfg.App
	})
	Context.AddConstructor(func(cfg Config) *config.PostgresConfig {
		return cfg.Pg
	})
	Context.AddConstructor(func(cfg Config) dbtool.Config {
		return cfg.Pg
	})
}
