package ipostgres

import (
	"database/sql"

	"github.com/typical-go/typical-go/appx"
)

// Connect to postgres database
func Connect(conf PGConfig) (*sql.DB, error) {
	return sql.Open("postgres", conf.ConnectionString())
}

// Create postgres infrastructure
func Create(config PGConfig) appx.DBInfra {
	return &DBInfra{
		config: config,
	}
}
