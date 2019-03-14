package repository

import (
	"log"

	"github.com/imantung/typical-go-server/config"
	"github.com/imantung/typical-go-server/db"
)

func init() {
	conf, _ := config.LoadConfigForTest()

	err := db.ResetTestDB(conf, "file://../../db/migrate")
	if err != nil {
		log.Fatal(err.Error())
	}
}
