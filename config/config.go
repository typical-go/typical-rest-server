package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typidb"
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

func GetDBToolConfig(cfg PostgresConfig) typidb.Config {
	return &cfg
}
