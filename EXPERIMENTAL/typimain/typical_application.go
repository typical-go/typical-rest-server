package typimain

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"gopkg.in/urfave/cli.v1"
)

// TypicalApplication represent typical application
type TypicalApplication struct {
	typictx.Context
}

// NewTypicalApplication return new instance of TypicalApplications
func NewTypicalApplication(context typictx.Context) *TypicalApplication {
	return &TypicalApplication{context}
}

// Cli return the command line interface
func (t *TypicalApplication) Cli() *cli.App {
	app := cli.NewApp()
	app.Name = t.Name
	app.Usage = ""
	app.Description = t.Description
	app.Version = t.Version
	app.Action = t.runActionFunc(t.ArcheType.GetAction())

	for _, cmd := range t.ArcheType.GetCommands() {
		app.Commands = append(app.Commands, cli.Command{
			Name:      cmd.Name,
			ShortName: cmd.ShortName,
			Action:    t.runActionFunc(cmd.Action),
		})
	}

	return app
}

func (t TypicalApplication) runActionFunc(action typictx.Action) interface{} {
	return func(ctx *cli.Context) error {
		return action.Start(typictx.ActionContext{
			CliContext: ctx,
			Context:    t.Context,
		})

	}
}
