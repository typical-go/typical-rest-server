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
	app.Action = t.invokeAction

	return app.Run(arguments)
}

func (t *TypicalApplication) invokeAction(ctx *cli.Context) error {
	return t.Container().Invoke(t.TypiApp.Action)
}
