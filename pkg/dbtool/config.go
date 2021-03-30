package dbtool

import (
	"database/sql"
	"fmt"
	"os"
)

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

var (
	DefaultEnvKeys = &EnvKeys{
		DBName: "DBNAME",
		DBUser: "DBUSER",
		DBPass: "DBPASS",
		Host:   "HOST",
		Port:   "PORT",
	}
)

func EnvKeysWithPrefix(prefix string) *EnvKeys {
	return &EnvKeys{
		DBName: prefix + "_" + DefaultEnvKeys.DBName,
		DBUser: prefix + "_" + DefaultEnvKeys.DBUser,
		DBPass: prefix + "_" + DefaultEnvKeys.DBPass,
		Host:   prefix + "_" + DefaultEnvKeys.Host,
		Port:   prefix + "_" + DefaultEnvKeys.Port,
	}
}

func openMySQL(c *Config) (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?tls=false&multiStatements=true",
		c.DBUser, c.DBPass, c.Host, c.Port, c.DBName,
	))
}

func openMySQLForAdmin(c *Config) (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf(
		"root:%s@tcp(%s:%s)/?tls=false&multiStatements=true",
		c.DBPass, c.Host, c.Port,
	))
}

func openPostgres(c *Config) (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPass, c.Host, c.Port, c.DBName,
	))
}

func openPostgresForAdmin(c *Config) (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/template1?sslmode=disable",
		c.DBUser, c.DBPass, c.Host, c.Port))
}

func (e *EnvKeys) GetConfig() *Config {
	return &Config{
		DBName: os.Getenv(e.DBName),
		DBUser: os.Getenv(e.DBUser),
		DBPass: os.Getenv(e.DBPass),
		Host:   os.Getenv(e.Host),
		Port:   os.Getenv(e.Port),
	}
}
