package typrails

import (
	"github.com/typical-go/typical-go/pkg/typbuild"
	"github.com/urfave/cli/v2"
)

// Rails of rails
type Rails struct{}

// New rails module
func New() *Rails {
	return &Rails{}
}

type rails struct { // TODO: remove this
	*typbuild.Context
}

// BuildCommands is commands to exectuce from Build-Tool
func (m *Rails) BuildCommands(c *typbuild.Context) []*cli.Command {
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
