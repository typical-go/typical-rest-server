package main

import (
	"log"
	"os"

	"github.com/typical-go/typical-rest-server/typical"
	"gopkg.in/urfave/cli.v1"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/lib/pq"
)

func main() {
	app := cli.NewApp()
	app.Name = typical.Context.Name
	app.Usage = ""
	app.Description = typical.Context.Description
	app.Version = typical.Context.Version
	app.Commands = Commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
