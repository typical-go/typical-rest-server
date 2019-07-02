package typicli

import (
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
	app.Name = t.context.Name + " (TYPICAL)"
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
		{Name: "update", ShortName: "u", Usage: "Update the typical binary", Action: updateTypical},
		{Name: "build", ShortName: "b", Usage: "Build the binary", Action: buildBinary},
		{Name: "test", Usage: "Run the Test", Action: runTest},
		{Name: "run", ShortName: "r", Usage: "Run the binary", Action: runBinary},
		{Name: "release", Usage: "Release the distribution", Action: releaseDistribution},
		{Name: "mock", Usage: "Generate mock class", Action: generateMock},
		{Name: "readme", Usage: "Generate readme", Action: t.generateReadme},

		{Name: "dep", Usage: "Vendoring the dependency", Action: vendoringDependency},
	}
}
