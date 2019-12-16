package typreadme

import (
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

// Module of readme
type Module struct{}

type readme struct {
	*typcore.Context
}

// BuildCommands to be shown in BuildTool
func (*Module) BuildCommands(c typcore.Cli) []*cli.Command {
	r := readme{Context: c.Context()}
	return []*cli.Command{
		r.generateCmd(),
	}
}
