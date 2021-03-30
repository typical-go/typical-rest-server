package typtool

import (
	"database/sql"
	"fmt"
	"os"
)

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

var (
	DefaultDBEnvKeys = &DBEnvKeys{
		DBName: "DBNAME",
		DBUser: "DBUSER",
		DBPass: "DBPASS",
		Host:   "HOST",
		Port:   "PORT",
	}
)

func DBEnvKeysWithPrefix(prefix string) *DBEnvKeys {
	return &DBEnvKeys{
		DBName: prefix + "_" + DefaultDBEnvKeys.DBName,
		DBUser: prefix + "_" + DefaultDBEnvKeys.DBUser,
		DBPass: prefix + "_" + DefaultDBEnvKeys.DBPass,
		Host:   prefix + "_" + DefaultDBEnvKeys.Host,
		Port:   prefix + "_" + DefaultDBEnvKeys.Port,
	}
}

func openMySQL(c *DBConfig) (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?tls=false&multiStatements=true",
		c.DBUser, c.DBPass, c.Host, c.Port, c.DBName,
	))
}

func openMySQLForAdmin(c *DBConfig) (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf(
		"root:%s@tcp(%s:%s)/?tls=false&multiStatements=true",
		c.DBPass, c.Host, c.Port,
	))
}

func openPostgres(c *DBConfig) (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPass, c.Host, c.Port, c.DBName,
	))
}

func openPostgresForAdmin(c *DBConfig) (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/template1?sslmode=disable",
		c.DBUser, c.DBPass, c.Host, c.Port))
}

func (e *DBEnvKeys) GetConfig() *DBConfig {
	return &DBConfig{
		DBName: os.Getenv(e.DBName),
		DBUser: os.Getenv(e.DBUser),
		DBPass: os.Getenv(e.DBPass),
		Host:   os.Getenv(e.Host),
		Port:   os.Getenv(e.Port),
	}
}
