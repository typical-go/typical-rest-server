package infra

import (
	"fmt"
)

// PostgresConfig postgres config
type PostgresConfig struct {
	DbName   string `envconfig:"PG_DBNAME" required:"true"`
	User     string `envconfig:"PG_USER" required:"true"`
	Password string `envconfig:"PG_PASSWORD" required:"true"`
	Host     string `envconfig:"PG_HOST" default:"localhost"`
	Port     int    `envconfig:"PG_PORT" default:"5432"`
}

// ConnectionString return connection string
func (p PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.DbName)
}

// ConnectionStringNoDB return connection string to template1 database
func (p PostgresConfig) ConnectionStringNoDB() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, "template1")
}
