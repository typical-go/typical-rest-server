package application

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/urfave/cli"
)

// Cli return the command line interface
func Cli(c *typictx.Context) *cli.App {
	app := cli.NewApp()
	app.Name = c.Name
	app.Usage = ""
	app.Description = c.Description
	app.Version = c.Version
	app.Action = typictx.ActionCommandFunction(c, c.Application)
	app.Before = func(ctx *cli.Context) error {
		return typienv.LoadEnvFile()
	}
	for _, cmd := range c.Application.Commands {
		app.Commands = append(app.Commands, cli.Command{
			Name:      cmd.Name,
			ShortName: cmd.ShortName,
			Usage:     cmd.Usage,
			Action:    cmd.ActionFunc.CommandFunction(c),
		})
	}
	return app
}
