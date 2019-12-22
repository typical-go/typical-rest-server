package typrails

import (
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

// Module of typrails
func Module() interface{} {
	return &module{}
}

type module struct{}

type rails struct {
	*typcore.Context
}

// BuildCommands is commands to exectuce from Build-Tool
func (m *module) BuildCommands(c *typcore.Context) []*cli.Command {
	r := rails{c}
	return []*cli.Command{
		{
			Name:  "rails",
			Usage: "Rails-like generation",
			Subcommands: []*cli.Command{
				r.scaffoldCmd(),
				r.repositoryCmd(),
			},
		},
	}
}
