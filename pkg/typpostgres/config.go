package typpostgres

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typcfg"
)

// Config is postgres configuration
type Config struct {
	DBName   string `required:"true"`
	User     string `required:"true" default:"postgres"`
	Password string `required:"true" default:"pgpass"`
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
}

// Configuration of Postgres
func Configuration() *typcfg.Configuration {
	return typcfg.NewConfiguration(DefaultConfigName, &Config{
		DBName:   DefaultUser,
		User:     DefaultUser,
		Password: DefaultPassword,
		Host:     DefaultHost,
		Port:     DefaultPort,
	})
}

// RetrieveConfig to retrieve postgres config from build-tool context
func RetrieveConfig(c *typbuildtool.Context) (cfg *Config, err error) {
	var v interface{}
	var ok bool

	if v, err = c.RetrieveConfig(DefaultConfigName); err != nil {
		return
	}

	if cfg, ok = v.(*Config); !ok {
		return nil, fmt.Errorf("Postgres: Get config for '%s' but invalid type", DefaultConfigName)
	}

	return
}
