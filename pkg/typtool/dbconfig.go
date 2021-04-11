package typtool

import (
	"os"
)

var DefaultDBEnvKeys = &DBEnvKeys{
	DBName: "DBNAME",
	DBUser: "DBUSER",
	DBPass: "DBPASS",
	Host:   "HOST",
	Port:   "PORT",
}

type (
	DBConfig struct {
		DBName string
		DBUser string
		DBPass string
		Host   string
		Port   string
	}
	DBEnvKeys DBConfig
)

//
// DBEnvKeys
//

func DBEnvKeysWithPrefix(prefix string) *DBEnvKeys {
	return &DBEnvKeys{
		DBName: prefix + "_" + DefaultDBEnvKeys.DBName,
		DBUser: prefix + "_" + DefaultDBEnvKeys.DBUser,
		DBPass: prefix + "_" + DefaultDBEnvKeys.DBPass,
		Host:   prefix + "_" + DefaultDBEnvKeys.Host,
		Port:   prefix + "_" + DefaultDBEnvKeys.Port,
	}
}

func (e *DBEnvKeys) Config() *DBConfig {
	return &DBConfig{
		DBName: os.Getenv(e.DBName),
		DBUser: os.Getenv(e.DBUser),
		DBPass: os.Getenv(e.DBPass),
		Host:   os.Getenv(e.Host),
		Port:   os.Getenv(e.Port),
	}
}
