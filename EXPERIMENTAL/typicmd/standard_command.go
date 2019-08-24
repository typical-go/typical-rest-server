package typicmd

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/urfave/cli"
)

// StandardCommands return standard commands for typical task tool
func StandardCommands(ctx *typictx.Context) []cli.Command {
	return []cli.Command{
		{
			Name:      "build",
			ShortName: "b",
			Usage:     "Build the binary",
			Action:    execCmd(ctx, BuildBinary),
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
			Action: execCmd(ctx, RunBinary),
		},
		{
			Name:      "docker",
			ShortName: "d",
			Usage:     "docker",
			Subcommands: []cli.Command{
				{
					Name:   "compose",
					Usage:  "Generate docker-compose.yaml",
					Action: execCmd(ctx, GenerateDockerCompose),
				},
				{
					Name:   "up",
					Usage:  "Create and start containers",
					Action: execCmd(ctx, DockerUp),
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
					Action: execCmd(ctx, DockerDown),
				},
			},
		},
		{
			Name:      "test",
			ShortName: "t",
			Usage:     "Run the testing",
			Action:    execCmd(ctx, RunTest),
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
			Action: execCmd(ctx, ReleaseDistribution),
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
			Action: execCmd(ctx, GenerateMock),
		},
		{
			Name:  "readme",
			Usage: "Generate readme document",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-commit",
					Usage: "Generate readme without auto commit",
				},
			},
			Action: execCmd(ctx, GenerateReadme),
		},
		{
			Name:        "clean",
			Usage:       "Clean project from generated file during build time",
			Description: "Remove binary folder, trigger `go clean --modcache`",
			Action:      execCmd(ctx, CleanProject),
		},
	}
}

func execCmd(ctx *typictx.Context, fn typictx.ActionFunc) interface{} {
	return func(cliCtx *cli.Context) error {
		return fn(&typictx.ActionContext{
			Context: ctx,
			Cli:     cliCtx,
		})
	}
}
