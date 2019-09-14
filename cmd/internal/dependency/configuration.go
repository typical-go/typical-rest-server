// Autogenerated by Typical-Go. DO NOT EDIT.

package dependency

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/typical-go/typical-rest-server/pkg/module/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/module/typredis"
	"github.com/typical-go/typical-rest-server/pkg/module/typserver"
	"github.com/typical-go/typical-rest-server/typical"
)

type Config struct {
	App    *config.Config
	Server *typserver.Config
	Pg     *typpostgres.Config
	Redis  *typredis.Config
}

func init() {
	typical.Context.AddConstructor(func() (*Config, error) {
		var cfg Config
		err := envconfig.Process("", &cfg)
		return &cfg, err
	})
	typical.Context.AddConstructor(func(cfg *Config) *config.Config {
		return cfg.App
	})
	typical.Context.AddConstructor(func(cfg *Config) *typserver.Config {
		return cfg.Server
	})
	typical.Context.AddConstructor(func(cfg *Config) *typpostgres.Config {
		return cfg.Pg
	})
	typical.Context.AddConstructor(func(cfg *Config) *typredis.Config {
		return cfg.Redis
	})
}
