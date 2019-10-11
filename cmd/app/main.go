package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/application"
	_ "github.com/typical-go/typical-rest-server/cmd/internal/dependency"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	app := application.NewApplication(typical.Context)
	err := app.Cli().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
