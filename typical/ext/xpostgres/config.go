package xpostgres

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// PGConfig postgres config
type PGConfig struct {
	DbName   string `required:"true"`
	User     string `required:"true"`
	Password string `required:"true"`
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
}

// ConnectionString return connection string
func (c PGConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DbName)
}

// ConnectionStringNoDB return connection string to template1 database
func (c PGConfig) ConnectionStringNoDB() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, "template1")
}

// LoadPostgresConfig load postgres configuration
func LoadPostgresConfig() (conf PGConfig, err error) {
	err = envconfig.Process("PG", &conf)
	return
}
