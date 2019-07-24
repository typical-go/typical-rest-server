package typitask

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"gopkg.in/urfave/cli.v1"
)

// TypicalTask contain typical command
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
			Action:    t.execCommand(BuildBinary),
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
			Action: t.execCommand(RunBinary),
		},
		{
			Name:      "test",
			ShortName: "t",
			Usage:     "Run the testing",
			Action:    t.execCommand(RunTest),
		},
		{
			Name:   "release",
			Usage:  "Release the distribution",
			Action: t.execCommand(ReleaseDistribution),
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
			Action: t.execCommand(GenerateMock),
		},
		{
			Name:   "readme",
			Usage:  "Generate readme document",
			Action: t.execCommand(GenerateReadme),
		},
		{
			Name:        "clean",
			Usage:       "Clean project from generated file during build time",
			Description: "Remove binary folder, trigger `go clean --modcache`",
			Action:      t.execCommand(CleanProject),
		},
		{
			Name:   "check",
			Usage:  "Checks all module status that required by the application",
			Action: t.execCommand(CheckStatus),
		},
	}
}

func (t *TypicalTask) execCommand(fn typictx.ActionFunc) interface{} {
	return func(cliCtx *cli.Context) error {
		return fn(typictx.ActionContext{
			Typical: t.Context,
			Cli:     cliCtx,
		})
	}
}
