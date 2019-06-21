package infra

import (
	"database/sql"

	"github.com/typical-go/typical-go/appx"
)

// ConnectPostgres connect to postgres database
func ConnectPostgres(conf PostgresConfig) (*sql.DB, error) {
	return sql.Open("postgres", conf.ConnectionString())
}

// CreatePostgresInfra create postgres infrastructure
func CreatePostgresInfra(config PostgresConfig) appx.DBInfra {
	return &PostgresInfra{
		config: config,
	}
}
