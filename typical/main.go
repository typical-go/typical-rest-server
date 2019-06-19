package main

import (
	"log"
	"os"

	"github.com/typical-go/typical-rest-server/typical/project"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = project.Ctx.Name
	app.Usage = ""
	app.Description = project.Ctx.Description
	app.Version = project.Ctx.Version
	app.Commands = Commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
