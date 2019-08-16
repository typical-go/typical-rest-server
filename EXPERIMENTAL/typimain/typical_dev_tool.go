package typimain

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typitask"
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
	app.Commands = t.StandardCommands()
	for key := range t.Modules {
		module := t.Modules[key]
		if module.Command != nil {
			app.Commands = append(app.Commands, typictx.ConvertToCLICommand(t.Context, module.Command))
		}
	}

	return app
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
					Usage: "Run the binary without build",
				},
			},
			Action: t.execCommand(typitask.RunBinary),
		},
		{
			Name:      "docker",
			ShortName: "d",
			Usage:     "docker",
			Subcommands: []cli.Command{
				{
					Name:   "compose",
					Usage:  "Generate docker-compose.yaml",
					Action: t.execCommand(typitask.GenerateDockerCompose),
				},
				{
					Name:   "up",
					Usage:  "Create and start containers",
					Action: t.execCommand(typitask.DockerUp),
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "no-gen",
							Usage: "Create and start containers without generate docker-compose.yaml",
						},
					},
				},
				{
					Name:   "down",
					Usage:  "Stop and remove containers, networks, images, and volumes",
					Action: t.execCommand(typitask.DockerDown),
				},
			},
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
				cli.BoolFlag{
					Name:  "no-test",
					Usage: "Release without run automated test",
				},
				cli.BoolFlag{
					Name:  "no-readme",
					Usage: "Release without generate readme",
				},
			},
			Action: t.execCommand(typitask.ReleaseDistribution),
		},
		{
			Name:  "mock",
			Usage: "Generate mock class",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-delete",
					Usage: "Generate mock class with delete previous generation",
				},
			},
			Action: t.execCommand(typitask.GenerateMock),
		},
		{
			Name:  "readme",
			Usage: "Generate readme document",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-auto-commit",
					Usage: "Generate readme without auto commit",
				},
			},
			Action: t.execCommand(typitask.GenerateReadme),
		},
		{
			Name:        "clean",
			Usage:       "Clean project from generated file during build time",
			Description: "Remove binary folder, trigger `go clean --modcache`",
			Action:      t.execCommand(typitask.CleanProject),
		},
	}
}

func (t *TypicalDevTool) execCommand(fn typictx.ActionFunc) interface{} {
	return func(cliCtx *cli.Context) error {
		return fn(&typictx.ActionContext{
			Context: t.Context,
			Cli:     cliCtx,
		})
	}
}
