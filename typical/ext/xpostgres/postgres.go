package xpostgres

import (
	"database/sql"

	"github.com/typical-go/typical-go/appx"
)

// Connect to postgres database
func Connect(conf Config) (*sql.DB, error) {
	return sql.Open("postgres", conf.ConnectionString())
}

// CreateDBInfra postgres infrastructure
func CreateDBInfra(config Config) appx.DBInfra {
	return &DBInfra{
		config: config,
	}
}
