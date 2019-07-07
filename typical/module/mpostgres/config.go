package mpostgres

import (
	"fmt"
)

// Config containr postgres configuration
type Config struct {
	DbName          string `required:"true" default:"default_dbname"`
	User            string `required:"true" default:"default_user"`
	Password        string `required:"true" default:"default_password"`
	Host            string `default:"localhost"`
	Port            int    `default:"5432"`
	MigrationSource string `default:"scripts/migration"`
}

// ConnectionString return connection string
func (c Config) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DbName)
}

// ConnectionStringNoDB return connection string to template1 database
func (c Config) ConnectionStringNoDB() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, "template1")
}
