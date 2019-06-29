package ipostgres

import (
	"database/sql"

	"github.com/typical-go/typical-go/appx"
)

// Connect to postgres database
func Connect(conf PGConfig) (*sql.DB, error) {
	return sql.Open("postgres", conf.ConnectionString())
}

// CreateDBInfra postgres infrastructure
func CreateDBInfra(config PGConfig) appx.DBInfra {
	return &DBInfra{
		config: config,
	}
}
