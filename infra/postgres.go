package infra

import "database/sql"

// ConnectPostgres connect to postgres database
func ConnectPostgres(conf PostgresConfig) (*sql.DB, error) {
	return sql.Open("postgres", conf.ConnectionString())
}

// CreatePostgresInfra create postgres infrastructure
func CreatePostgresInfra(config PostgresConfig) DBInfra {
	return &PostgresInfra{
		config: config,
	}
}
