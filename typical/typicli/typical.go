package typicli

import (
	"fmt"

	"github.com/typical-go/typical-rest-server/typical/appctx"
	"gopkg.in/urfave/cli.v1"
)

// Typical program
type Typical struct {
	context *appctx.Context
}

// NewTypical return new instance of Typical
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

	for key := range t.context.Commands {
		command := t.context.Commands[key]
		app.Commands = append(app.Commands, command)
	}
	return app.Run(arguments)
}

func (t *Typical) standardTypicalCommand() []cli.Command {
	return []cli.Command{
		{Name: "build", Usage: "Build the binary", Action: notImplemented},
		{Name: "test", Usage: "Run the Test", Action: notImplemented},
		{Name: "run", Usage: "Run the binary", Action: notImplemented},
		{Name: "dep", Usage: "Vendoring the dependency", Action: notImplemented},
		{Name: "release", Usage: "Release the distribution", Action: notImplemented},
		{Name: "mock", Usage: "Generate mock class", Action: mock},
		{Name: "readme", Usage: "Generate readme", Action: t.generateReadme},
	}
}

// NotImplement for not implemented function
func notImplemented(ctx *cli.Context) {
	fmt.Println("Not implemented")
}
