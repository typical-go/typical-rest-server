package db

import (
	"fmt"

	"github.com/imantung/typical-go-server/config"
)

func connectionString(conf config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.DbUser, conf.DbPassword, conf.DbHost, conf.DbPort, conf.DbName)
}

func connectionStringWithDBName(conf config.Config, dbname string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.DbUser, conf.DbPassword, conf.DbHost, conf.DbPort, dbname)
}
