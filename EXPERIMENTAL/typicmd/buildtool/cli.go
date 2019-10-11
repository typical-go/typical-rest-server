package buildtool

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/urfave/cli"
)

// Cli return the command line interface
func Cli(c *typictx.Context) *cli.App {
	app := cli.NewApp()
	app.Name = c.Name
	app.Usage = ""
	app.Description = c.Description
	app.Version = c.Version
	app.Before = func(ctx *cli.Context) error {
		return c.Validate()
	}
	for _, cmd := range commands(c) {
		app.Commands = append(app.Commands, cmd.CliCommand(c))
	}
	return app
}
