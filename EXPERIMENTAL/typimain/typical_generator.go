package typimain

import (
	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typigen"
	"gopkg.in/urfave/cli.v1"
)

// TypicalGenerator represent typical generator
type TypicalGenerator struct {
	typictx.Context
}

// NewTypicalGenerator return new instance of TypicalCli
func NewTypicalGenerator(context typictx.Context) *TypicalGenerator {
	return &TypicalGenerator{
		Context: context,
	}
}

// Cli return the command line interface
func (t *TypicalGenerator) Cli() *cli.App {
	app := cli.NewApp()
	app.Action = func(ctx *cli.Context) error {
		return runn.Execute(
			typienv.WriteEnvIfNotExist(t.Context),
			typienv.LoadEnv(),
			typigen.MainAppGenerated(t.Context),
			typigen.MainDevToolGenerated(t.Context),
			typigen.TypicalGenerated(t.Context),
		)
	}

	return app
}
