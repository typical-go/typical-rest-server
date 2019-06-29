package typicli

import (
	"fmt"

	"github.com/typical-go/typical-rest-server/typical/appctx"
	"gopkg.in/urfave/cli.v1"
)

type Typical struct {
	context *appctx.Context
}

func NewTypical(context *appctx.Context) *Typical {
	return &Typical{
		context: context,
	}
}

// Run the typical task cli
func (t *Typical) Run(arguments []string) error {
	app := cli.NewApp()
	app.Name = t.context.Name
	app.Usage = ""
	app.Description = t.context.Description
	app.Version = t.context.Version

	app.Commands = t.standardTypicalCommand()
	for key := range t.context.Modules {
		module := t.context.Modules[key]
		app.Commands = append(app.Commands, module.Command())
	}
	return app.Run(arguments)
}

func (t *Typical) standardTypicalCommand() []cli.Command {
	return []cli.Command{
		{Name: "build", Usage: "Build the binary", Action: notImplemented},
		{Name: "test", Usage: "Run the Test", Action: notImplemented},
		{Name: "run", Usage: "Run the binary", Action: notImplemented},
		{Name: "vendoring", Usage: "Vendoring the dependency", Action: notImplemented},
		{Name: "release", Usage: "Release the distribution", Action: notImplemented},
		{Name: "mock", Usage: "Generate mock class", Action: mock},
		{Name: "readme", Usage: "Generate readme", Action: t.generateReadme},
	}
}

// NotImplement for not implemented function
func notImplemented(ctx *cli.Context) {
	fmt.Println("Not implemented")
}
