// Package config contain configuration of the project.
// It's a good idea for configuration to have its own package to prevent the circular import because it naturally will be used by multiple packages.
// Don't forget to create test for each configuration field to make sure no typo and get the expected value/
package config

import "github.com/kelseyhightower/envconfig"

// Config type
type Config struct {
	Address string `envconfig:"ADDRESS"`
}

// Load the configuration
func Load() (conf Config, err error) {
	err = envconfig.Process(Prefix, &conf)

	return
}
