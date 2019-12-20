package typrails

import (
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

// Module of typrails
type Module struct{}

// BuildCommands is commands to exectuce from Build-Tool
func (m *Module) BuildCommands(c *typcore.Context) []*cli.Command {
	r := rails{c}
	return []*cli.Command{
		r.scaffoldCmd(),
	}
}
