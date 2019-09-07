package typimain

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typigen"
	"github.com/urfave/cli"
)

// TypicalGenerator represent typical generator
type TypicalGenerator struct {
	*typictx.Context
}

// NewTypicalGenerator return new instance of TypicalCli
func NewTypicalGenerator(context *typictx.Context) *TypicalGenerator {
	return &TypicalGenerator{
		Context: context,
	}
}

// Cli return the command line interface
func (g *TypicalGenerator) Cli() *cli.App {
	app := cli.NewApp()
	app.Action = g.run
	return app
}

func (g *TypicalGenerator) run(ctx *cli.Context) error {
	return typigen.Generate(g.Context)
}
