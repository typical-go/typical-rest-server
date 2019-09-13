package typimain

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/urfave/cli"
)

// TypicalBuildTool represent typical task tool application
type TypicalBuildTool struct {
	*typictx.Context
}

// NewTypicalBuildTool return new instance of TypicalCli
func NewTypicalBuildTool(context *typictx.Context) *TypicalBuildTool {
	return &TypicalBuildTool{
		Context: context,
	}
}

// Cli return the command line interface
func (t *TypicalBuildTool) Cli() *cli.App {
	app := cli.NewApp()
	app.Name = t.Name + " (TYPICAL)"
	app.Usage = ""
	app.Description = t.Description
	app.Version = t.Version
	for _, cmd := range t.commands() {
		app.Commands = append(app.Commands, cmd.CliCommand(t.Context))
	}
	return app
}

func (t *TypicalBuildTool) commands() []*typictx.Command {
	cmds := typicmd.StandardCommands(t.Context)
	for _, module := range t.Modules {
		if module.Command != nil {
			cmds = append(cmds, module.Command)
		}
	}
	return cmds
}
