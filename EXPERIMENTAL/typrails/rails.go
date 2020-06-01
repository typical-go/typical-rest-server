package typrails

import (
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/urfave/cli/v2"
)

// Rails of rails
type Rails struct{}

// New rails module
func New() *Rails {
	return &Rails{}
}

// Commands to exectuce from Build-Tool
func (m *Rails) Commands(c *typgo.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "rails",
			Usage: "Rails-like generation",
			Subcommands: []*cli.Command{
				scaffoldCmd(c),
				repositoryCmd(c),
			},
		},
	}
}
