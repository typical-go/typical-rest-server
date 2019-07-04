package xpostgres

import (
	"database/sql"
)

// Connect to postgres database
func Connect(conf Config) (*sql.DB, error) {
	return sql.Open("postgres", conf.ConnectionString())
}
