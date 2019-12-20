package typreadme

import (
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

// Module of readme
func Module() interface{} {
	return &module{}
}

type module struct{}

type readme struct {
	*typcore.Context
}

// BuildCommands to be shown in BuildTool
func (*module) BuildCommands(c *typcore.Context) []*cli.Command {
	r := readme{Context: c}
	return []*cli.Command{
		r.generateCmd(),
	}
}
