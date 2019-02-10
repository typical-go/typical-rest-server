package db

import (
	"database/sql"

	"github.com/imantung/go-helper/dbkit"
	"github.com/imantung/typical-go-server/config"

	// load the driver
	_ "github.com/lib/pq"
)

// Connect to database
func Connect(conf config.Config) (*sql.DB, error) {
	pgConf := dbkit.PgConfig{
		Host:     conf.DbHost,
		Port:     conf.DbPort,
		DbName:   conf.DbName,
		User:     conf.DbUser,
		Password: conf.DbPassword,
		SslMode:  "disable",
	}

	return sql.Open("postgres", pgConf.ConnectionString())
}
