package main

import (
	"log"
	"os"

	"github.com/typical-go/typical-rest-server/typical/projctx"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = projctx.Name()
	app.Usage = ""
	app.Description = projctx.Description()
	app.Version = projctx.Version()
	app.Commands = Commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
