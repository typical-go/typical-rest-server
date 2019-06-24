package ipostgres

import (
	"fmt"
)

// PGConfig postgres config
type PGConfig struct {
	DbName   string `envconfig:"PG_DBNAME" required:"true"`
	User     string `envconfig:"PG_USER" required:"true"`
	Password string `envconfig:"PG_PASSWORD" required:"true"`
	Host     string `envconfig:"PG_HOST" default:"localhost"`
	Port     int    `envconfig:"PG_PORT" default:"5432"`
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
