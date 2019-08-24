package typimain

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/urfave/cli"
)

// TypicalDevTool represent typical task tool application
type TypicalDevTool struct {
	*typictx.Context
}

// NewTypicalDevTool return new instance of TypicalCli
func NewTypicalDevTool(context *typictx.Context) *TypicalDevTool {
	return &TypicalDevTool{
		Context: context,
	}
}

// Cli return the command line interface
func (t *TypicalDevTool) Cli() *cli.App {
	app := cli.NewApp()
	app.Name = t.Name + " (TYPICAL)"
	app.Usage = ""
	app.Description = t.Description
	app.Version = t.Version
	app.Commands = typicmd.StandardCommands(t.Context)
	for key := range t.Modules {
		module := t.Modules[key]
		if module.Command != nil {
			app.Commands = append(app.Commands, typictx.ConvertToCLICommand(t.Context, module.Command))
		}
	}
	return app
}
