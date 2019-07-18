package typimain

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typitask"
	"gopkg.in/urfave/cli.v1"
)

// TypicalDevTool represent typical task tool application
type TypicalDevTool struct {
	typitask.TypicalTask
}

// NewTypicalTaskTool return new instance of TypicalCli
func NewTypicalTaskTool(context typictx.Context) *TypicalDevTool {
	return &TypicalDevTool{
		typitask.TypicalTask{
			Context: context,
		},
	}
}

// Cli return the command line interface
func (t *TypicalDevTool) Cli() *cli.App {
	app := cli.NewApp()
	app.Name = t.Name + " (TYPICAL)"
	app.Usage = ""
	app.Description = t.Description
	app.Version = t.Version
	app.Commands = t.StandardCommands()
	for key := range t.Modules {
		module := t.Modules[key]

		if len(module.Commands) > 0 {
			app.Commands = append(app.Commands, cli.Command{
				Name:        module.Name,
				ShortName:   module.ShortName,
				Usage:       module.Usage,
				Subcommands: module.Commands,
			})
		}
	}

	for key := range t.Commands {
		command := t.Commands[key]
		app.Commands = append(app.Commands, command)
	}

	// TODO: put it at before run of cli.
	typienv.ExportProjectEnv()

	return app
}
