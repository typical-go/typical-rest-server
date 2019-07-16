package typitask

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"gopkg.in/urfave/cli.v1"
)

type TypicalTask struct {
	typictx.Context
}

// StandardCommands return standard commands for typical task tool
func (t *TypicalTask) StandardCommands() []cli.Command {
	return []cli.Command{
		{
			Name:      "build",
			ShortName: "b",
			Usage:     "Build the binary",
			Action:    t.buildBinary,
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
			Action: t.runBinary,
		},
		{
			Name:   "release",
			Usage:  "Release the distribution",
			Action: t.releaseDistribution,
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
			Action: t.generateMock,
		},
		{
			Name:   "readme",
			Usage:  "Generate readme document",
			Action: t.generateReadme,
		},
		{
			Name:        "clean",
			Usage:       "Clean project from generated file during build time",
			Description: "Remove binary folder, trigger `go clean --modcache`",
			Action:      t.cleanProject,
		},
	}
}
