package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/dbtool"
)

// TODO: to be generated

type Config struct {
	App AppConfig
	Pg  PostgresConfig
}

func LoadConfig() (cfg Config, err error) {
	err = envconfig.Process("", &cfg)
	return
}

func GetAppConfig(cfg Config) AppConfig {
	return cfg.App
}

func GetPgConfig(cfg Config) PostgresConfig {
	return cfg.Pg
}

func GetDBToolConfig(cfg PostgresConfig) dbtool.Config {
	return &cfg
}
