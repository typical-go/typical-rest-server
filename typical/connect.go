package typical

import (
	"database/sql"

	"github.com/typical-go/typical-rest-server/config"
)

func connectDB(conf config.Config) (*sql.DB, error) {
	return conf.Postgres.Open()
}
