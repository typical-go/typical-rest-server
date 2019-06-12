package main

import (
	"log"
	"os"

	"github.com/typical-go/typical-rest-server/typical/projctx"
	"github.com/urfave/cli"
)

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = projctx.Name()
	cliApp.Usage = projctx.Usage()
	cliApp.UsageText = projctx.Example()
	cliApp.Version = projctx.Version()
	cliApp.Commands = Commands()

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
