// Package config contain configuration of the project.
// It's a good idea for configuration to have its own package to prevent the circular import because it naturally will be used by multiple packages.
// Don't forget to create test for each configuration field to make sure no typo and get the expected value/
package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/app/server"
	"github.com/typical-go/typical-rest-server/infra"
)

var Prefix = "APP"

// Config type
// Check https://github.com/kelseyhightower/envconfig#struct-tag-support for more help
type Config struct {
	server.Config
	infra.PostgresConfig
}

// LoadConfig return new instance of config
func LoadConfig() (config server.Config, err error) {
	err = envconfig.Process(Prefix, &config)
	return
}

// LoadPostgresConfig load postgres configuration
func LoadPostgresConfig() (conf infra.PostgresConfig, err error) {
	err = envconfig.Process(Prefix, &conf)
	return
}
