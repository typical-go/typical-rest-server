package typdb

import (
	"os"
)

var DefaultEnvKeys = &EnvKeys{
	DBName: "DBNAME",
	DBUser: "DBUSER",
	DBPass: "DBPASS",
	Host:   "HOST",
	Port:   "PORT",
}

type (
	Config struct {
		DBName string
		DBUser string
		DBPass string
		Host   string
		Port   string
	}
	EnvKeys Config
)

//
// EnvKeys
//

func EnvKeysWithPrefix(prefix string) *EnvKeys {
	return &EnvKeys{
		DBName: prefix + "_" + DefaultEnvKeys.DBName,
		DBUser: prefix + "_" + DefaultEnvKeys.DBUser,
		DBPass: prefix + "_" + DefaultEnvKeys.DBPass,
		Host:   prefix + "_" + DefaultEnvKeys.Host,
		Port:   prefix + "_" + DefaultEnvKeys.Port,
	}
}

func (e *EnvKeys) Config() *Config {
	return &Config{
		DBName: os.Getenv(e.DBName),
		DBUser: os.Getenv(e.DBUser),
		DBPass: os.Getenv(e.DBPass),
		Host:   os.Getenv(e.Host),
		Port:   os.Getenv(e.Port),
	}
}
