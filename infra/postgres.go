package infra

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/database/postgres"
)

// Postgres infrastructure model
type Postgres struct {
	DbName   string `envconfig:"PG_DBNAME" required:"true"`
	User     string `envconfig:"PG_USER" required:"true"`
	Password string `envconfig:"PG_PASSWORD" required:"true"`
	Host     string `envconfig:"PG_HOST" default:"localhost"`
	Port     int    `envconfig:"PG_PORT" default:"5432"`
}

// Open connection to postgres
func (p Postgres) Open() (*sql.DB, error) {
	return sql.Open("postgres", p.ConnectionString())
}

// ConnectionString return connection string
func (p Postgres) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.DbName)
}

// ConnectionStringTemplate1 return connection string to template1 database
func (p Postgres) ConnectionStringTemplate1() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, "template1")
}
