package application

import (
	log "github.com/sirupsen/logrus"

	"os"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/urfave/cli"
)

// Run the application
func Run(c *typictx.Context) {
	app := cli.NewApp()
	app.Name = c.Name
	app.Usage = ""
	app.Description = c.Description
	app.Version = c.Version
	app.Action = action(c, c.StartFunc)
	app.Before = func(ctx *cli.Context) error {
		return typienv.LoadEnvFile()
	}
	for _, cmd := range c.Application.Commands {
		app.Commands = append(app.Commands, cli.Command{
			Name:      cmd.Name,
			ShortName: cmd.ShortName,
			Usage:     cmd.Usage,
			Action:    action(c, cmd.ActionFunc),
		})
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// ActionCommandFunction to get command function fo action
func action(ctx *typictx.Context, action interface{}) interface{} {
	return runner{
		Context: ctx,
		action:  action,
	}.Run
}
