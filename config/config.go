// Package config contain configuration of the project.
// It's a good idea for configuration to have its own package to prevent the circular import because it naturally will be used by multiple packages.
// Don't forget to create test for each configuration field to make sure no typo and get the expected value/
package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config type
// Check https://github.com/kelseyhightower/envconfig#struct-tag-support for more help
type Config struct {
	Address string `envconfig:"ADDRESS" required:"true"`

	DbName     string `envconfig:"DB_NAME" required:"true"`
	DbUser     string `envconfig:"DB_USER" required:"true"`
	DbPassword string `envconfig:"DB_PASSWORD" required:"true"`
	DbHost     string `envconfig:"DB_HOST" default:"localhost"`
	DbPort     int    `envconfig:"DB_PORT" default:"5432"`
}

// LoadConfig return new instance of config
func LoadConfig() (conf Config, err error) {
	err = envconfig.Process(Prefix, &conf)
	return
}

// LoadConfigForTest return config for testing environment
func LoadConfigForTest() (conf Config, err error) {
	conf, err = LoadConfig()
	conf.DbName = conf.DbName + "_test"
	return
}

// Details return details of configuration
func Details() []ConfigDetail {
	return details(&Config{})
}
