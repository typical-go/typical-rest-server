package db

import (
	"fmt"

	"github.com/imantung/typical-go-server/config"
)

func connectionString(conf config.Config) string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		conf.DbHost, conf.DbPort, conf.DbName, conf.DbUser, conf.DbPassword)
}

func connectionStringWithDBName(conf config.Config, dbname string) string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		conf.DbHost, conf.DbPort, dbname, conf.DbUser, conf.DbPassword)
}
