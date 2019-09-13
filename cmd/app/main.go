package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	app := typicmd.NewTypicalApplication(typical.Context)
	err := app.Cli().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
