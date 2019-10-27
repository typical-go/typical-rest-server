package buildtool

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/urfave/cli"
)

// Run the build tool
func Run(c *typictx.Context) {
	buildtool := buildtool{Context: c}
	app := cli.NewApp()
	app.Name = c.Name
	app.Usage = ""
	app.Description = c.Description
	app.Version = c.Version
	app.Before = func(ctx *cli.Context) error {
		return c.Preparing()
	}
	for _, cmd := range buildtool.commands() {
		app.Commands = append(app.Commands, cmd.CliCommand(c))
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
