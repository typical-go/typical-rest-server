package application

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/urfave/cli"
)

// Application represent typical application
type Application struct {
	*typictx.Context
}

// NewApplication return new instance of TypicalApplications
func NewApplication(context *typictx.Context) *Application {
	return &Application{context}
}

// Cli return the command line interface
func (t *Application) Cli() *cli.App {
	app := cli.NewApp()
	app.Name = t.Name
	app.Usage = ""
	app.Description = t.Description
	app.Version = t.Version
	app.Action = typictx.ActionCommandFunction(t.Context, t.Application)
	app.Before = func(ctx *cli.Context) error {
		return typienv.LoadEnv()
	}

	for _, cmd := range t.Application.Commands {
		app.Commands = append(app.Commands, cli.Command{
			Name:      cmd.Name,
			ShortName: cmd.ShortName,
			Usage:     cmd.Usage,
			Action:    cmd.ActionFunc.CommandFunction(t.Context),
		})
	}
	return app
}
