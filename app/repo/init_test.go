package repo

import (
	"log"

	"github.com/typical-go/typical-rest-server/config"
	"github.com/typical-go/typical-rest-server/db"
)

func init() {
	conf, _ := config.LoadConfig()

	err := db.ResetTestDB(conf, "file://../../db/migrate")
	if err != nil {
		log.Fatal(err.Error())
	}
}
