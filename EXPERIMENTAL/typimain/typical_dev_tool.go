package typimain

import (
	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typigen"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typitask"
	"gopkg.in/urfave/cli.v1"
)

// TypicalDevTool represent typical task tool application
type TypicalDevTool struct {
	typictx.Context
}

// NewTypicalDevTool return new instance of TypicalCli
func NewTypicalDevTool(context typictx.Context) *TypicalDevTool {
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
	app.Before = t.beforeAction
	app.Commands = t.StandardCommands()
	for key := range t.Modules {
		module := t.Modules[key]

		if len(module.Commands) > 0 {
			command := cli.Command{
				Name:      module.Name,
				ShortName: module.ShortName,
				Usage:     module.Usage,
			}
			for i := range module.Commands {
				subCommand := module.Commands[i]
				command.Subcommands = append(command.Subcommands, cli.Command{
					Name:      subCommand.Name,
					ShortName: subCommand.ShortName,
					Usage:     subCommand.Usage,
					Action:    runActionFunc(t.Context, subCommand.ActionFunc),
				})
			}
			app.Commands = append(app.Commands, command)
		}
	}

	return app
}

func (t *TypicalDevTool) beforeAction(cliCtx *cli.Context) (err error) {
	return runn.Execute(
		typienv.WriteEnvIfNotExist(t.Context),
		typigen.MainAppGenerated(t.Context),
		typigen.MainDevToolGenerated(t.Context),
		typigen.TypicalGenerated(t.Context),
	)
}

// StandardCommands return standard commands for typical task tool
func (t *TypicalDevTool) StandardCommands() []cli.Command {
	return []cli.Command{
		{
			Name:      "build",
			ShortName: "b",
			Usage:     "Build the binary",
			Action:    t.execCommand(typitask.BuildBinary),
		},
		{
			Name:      "run",
			ShortName: "r",
			Usage:     "Run the binary",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-build",
					Usage: "Run the binary",
				},
			},
			Action: t.execCommand(typitask.RunBinary),
		},
		{
			Name:      "test",
			ShortName: "t",
			Usage:     "Run the testing",
			Action:    t.execCommand(typitask.RunTest),
		},
		{
			Name:  "release",
			Usage: "Release the distribution",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "github-token",
					Usage: "Release to github using github-token",
				},
				cli.BoolFlag{
					Name:  "alpha",
					Usage: "Release as alpha version (pre-release)",
				},
				cli.BoolFlag{
					Name:  "no-test",
					Usage: "Release without run automated test",
				},
				cli.BoolFlag{
					Name:  "no-readme",
					Usage: "Release without generate readme",
				},
				cli.BoolFlag{
					Name:  "force",
					Usage: "Release even if it already exist",
				},
				cli.StringFlag{
					Name:  "note",
					Usage: "Note for this release",
				},
			},
			Action: t.execCommand(typitask.ReleaseDistribution),
		},
		{
			Name:  "mock",
			Usage: "Generate mock class",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "new",
					Usage: "Clean the mock package as new generation",
				},
			},
			Action: t.execCommand(typitask.GenerateMock),
		},
		{
			Name:   "readme",
			Usage:  "Generate readme document",
			Action: t.execCommand(typitask.GenerateReadme),
		},
		{
			Name:        "clean",
			Usage:       "Clean project from generated file during build time",
			Description: "Remove binary folder, trigger `go clean --modcache`",
			Action:      t.execCommand(typitask.CleanProject),
		},
		{
			Name:   "check",
			Usage:  "Checks all module status that required by the application",
			Action: t.execCommand(typitask.CheckStatus),
		},
	}
}

func (t *TypicalDevTool) execCommand(fn typictx.ActionFunc) interface{} {
	return func(cliCtx *cli.Context) error {
		return fn(typictx.ActionContext{
			Typical: t.Context,
			Cli:     cliCtx,
		})
	}
}
