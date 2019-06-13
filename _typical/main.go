package main

import (
	"log"
	"os"

	"github.com/typical-go/typical-rest-server/_typical/project"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = project.Context.Name
	app.Usage = ""
	app.Description = project.Context.Description
	app.Version = project.Context.Version
	app.Commands = Commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
