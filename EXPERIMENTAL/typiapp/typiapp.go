package typiapp

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

// Run the typical application
func (t *TypicalApplication) Run(arguments []string) error {
	app := cli.NewApp()
	app.Name = t.Name
	app.Usage = ""
	app.Description = t.Description
	app.Version = t.Version
	app.Action = t.action

	for _, cmd := range t.ArcheType.GetCommands() {
		app.Commands = append(app.Commands, cli.Command{
			Name:      cmd.Name,
			ShortName: cmd.ShortName,
			Action:    t.invoke(cmd.Action),
		})
	}

	return app.Run(arguments)
}

func (t *TypicalApplication) action(ctx *cli.Context) {
	t.ArcheType.StartApplication(typictx.StartContext{
		CliContext: ctx,
		Context:    t.Context,
	})
}

func (t *TypicalApplication) invoke(f interface{}) interface{} {
	return func(ctx *cli.Context) error {
		return t.Container().Invoke(f)
	}
}
