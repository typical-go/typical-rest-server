package typical

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/app/server"
	"github.com/typical-go/typical-rest-server/infra"
)

// AllConfig all configuration
type AllConfig struct {
	server.Config
	infra.PostgresConfig
}

// LoadConfig return new instance of config
func LoadConfig() (config server.Config, err error) {
	err = envconfig.Process(Context.ConfigPrefix, &config)
	return
}

// LoadPostgresConfig load postgres configuration
func LoadPostgresConfig() (conf infra.PostgresConfig, err error) {
	err = envconfig.Process(Context.ConfigPrefix, &conf)
	return
}
