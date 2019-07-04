package typiapp

import (
	"github.com/typical-go/typical-rest-server/typical/appctx"
	"gopkg.in/urfave/cli.v1"
)

// TypicalApplication represent typical application
type TypicalApplication struct {
	appctx.Context
}

// NewTypicalApplication return new instance of TypicalApplications
func NewTypicalApplication(context appctx.Context) *TypicalApplication {
	return &TypicalApplication{context}
}

// Run the typical application
func (t *TypicalApplication) Run(arguments []string) error {
	app := cli.NewApp()
	app.Name = t.Name
	app.Usage = ""
	app.Description = t.Description
	app.Version = t.Version
	app.Action = t.invoke(t.TypiApp.Action)

	for i := range t.TypiApp.Commands {
		cmd := t.TypiApp.Commands[i]
		app.Commands = append(app.Commands, cli.Command{
			Name:      cmd.Name,
			ShortName: cmd.ShortName,
			Action:    t.invoke(cmd.Action),
		})
	}

	return app.Run(arguments)
}

func (t *TypicalApplication) invoke(f interface{}) interface{} {
	return func(ctx *cli.Context) error {
		return t.Container().Invoke(f)
	}
}
