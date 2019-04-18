package app

import (
	"database/sql"

	"github.com/imantung/typical-go-server/config"
)

func connectDB(conf config.Config) (*sql.DB, error) {
	return conf.Postgres.Open()
}
