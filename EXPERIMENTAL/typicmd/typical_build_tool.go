package typicmd

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/buildtool"
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
	app.Before = func(ctx *cli.Context) error {
		return t.Validate()
	}
	for _, cmd := range t.commands() {
		app.Commands = append(app.Commands, cmd.CliCommand(t.Context))
	}
	return app
}

func (t *TypicalBuildTool) commands() []*typictx.Command {
	cmds := buildtool.StandardCommands(t.Context)
	for _, module := range t.Modules {
		if module.Command != nil {
			cmds = append(cmds, module.Command)
		}
	}
	return cmds
}
