package buildtool

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/urfave/cli"
)

func commands(c *typictx.Context) (cmds []*typictx.Command) {
	cmds = []*typictx.Command{
		{
			Name:       "build",
			ShortName:  "b",
			Usage:      "Build the binary",
			ActionFunc: buildBinary,
		},
		{
			Name:       "clean",
			ShortName:  "c",
			Usage:      "Clean the project from generated file during build time",
			ActionFunc: cleanProject,
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
			ActionFunc: runBinary,
		},
		{
			Name:       "test",
			ShortName:  "t",
			Usage:      "Run the testing",
			ActionFunc: runTesting,
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
					Name:  "no-github",
					Usage: "Release without create github release",
				},
				cli.BoolFlag{
					Name:  "force",
					Usage: "Release by passed all validation",
				},
				cli.BoolFlag{
					Name:  "alpha",
					Usage: "Release for alpha version",
				},
			},
			ActionFunc: releaseDistribution,
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
			ActionFunc: generateMock,
		},
		{
			Name:       "readme",
			Usage:      "Generate readme document",
			ActionFunc: generateReadme,
		},
		{
			Name:       "docker",
			Usage:      "Docker utility",
			BeforeFunc: typienv.LoadEnvFile,
			SubCommands: []*typictx.Command{
				{
					Name:       "compose",
					Usage:      "Generate docker-compose.yaml",
					ActionFunc: dockerCompose,
				},
				{
					Name:  "up",
					Usage: "Create and start containers",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "no-compose",
							Usage: "Create and start containers without generate docker-compose.yaml",
						},
					},
					ActionFunc: dockerUp,
				},
				{
					Name:       "down",
					Usage:      "Stop and remove containers, networks, images, and volumes",
					ActionFunc: dockerDown,
				},
			},
		},
	}
	for _, module := range c.Modules {
		if module.Command != nil {
			cmds = append(cmds, module.Command)
		}
	}
	return
}
