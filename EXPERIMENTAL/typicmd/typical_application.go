package typicmd

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/urfave/cli"
)

// TypicalApplication represent typical application
type TypicalApplication struct {
	*typictx.Context
}

// NewTypicalApplication return new instance of TypicalApplications
func NewTypicalApplication(context *typictx.Context) *TypicalApplication {
	return &TypicalApplication{context}
}

// Cli return the command line interface
func (t *TypicalApplication) Cli() *cli.App {
	app := cli.NewApp()
	app.Name = t.Name
	app.Usage = ""
	app.Description = t.Description
	app.Version = t.Version
	app.Action = typictx.ActionCommandFunction(t.Context, t.Application)
	app.Before = t.beforeApplication

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

func (t *TypicalApplication) beforeApplication(ctx *cli.Context) error {
	return typienv.LoadEnv()
}
