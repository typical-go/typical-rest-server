package typrails

import (
	"github.com/typical-go/typical-go/pkg/typbuild"
	"github.com/urfave/cli/v2"
)

// New rails module
func New() *Module {
	return &Module{}
}

// Module of rails
type Module struct{}

type rails struct {
	*typbuild.Context
}

// BuildCommands is commands to exectuce from Build-Tool
func (m *Module) BuildCommands(c *typbuild.Context) []*cli.Command {
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
