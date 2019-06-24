package typical

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/typical/infra/ipostgres"
)

// AllConfig all configuration
type AllConfig struct {
	app.Config
	ipostgres.PGConfig
}

// LoadConfig return new instance of config
func LoadConfig() (config app.Config, err error) {
	err = envconfig.Process(Context.ConfigPrefix, &config)
	return
}

// LoadPostgresConfig load postgres configuration
func LoadPostgresConfig() (conf ipostgres.PGConfig, err error) {
	err = envconfig.Process(Context.ConfigPrefix, &conf)
	return
}
