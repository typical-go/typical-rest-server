package main

import (
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/lib/pq"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	err := typical.Context.Run(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
