package typimain

import (
	"github.com/typical-go/typical-rest-server/experimental/typictx"
	"gopkg.in/urfave/cli.v1"
)

// TypicalTask represent typical CLI
type TypicalTask struct {
	typictx.Context
}

// NewTypicalCli return new instance of TypicalCli
func NewTypicalCli(context typictx.Context) *TypicalTask {
	return &TypicalTask{context}
}

// Run the typical task cli
func (t *TypicalTask) Run(arguments []string) error {
	app := cli.NewApp()
	app.Name = t.Name + " (TYPICAL)"
	app.Usage = ""
	app.Description = t.Description
	app.Version = t.Version
	app.Commands = t.standardTypicalCommand()
	for key := range t.Modules {
		module := t.Modules[key]
		if module.Command != nil {
			app.Commands = append(app.Commands, *module.Command)
		}

	}

	for key := range t.Commands {
		command := t.Commands[key]
		app.Commands = append(app.Commands, command)
	}

	// NOTE: export the enviroment before run
	exportEnviroment()

	return app.Run(arguments)
}

func (t *TypicalTask) standardTypicalCommand() []cli.Command {
	return []cli.Command{
		{
			Name:      "update",
			ShortName: "u",
			Usage:     "Update the typical",
			Action:    t.updateTypical,
		},
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
			Name:   "mock",
			Usage:  "Generate mock class",
			Action: t.generateMock,
		},
		{
			Name:   "readme",
			Usage:  "Generate readme document",
			Action: t.generateReadme,
		},
	}
}
