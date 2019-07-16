package typictx

import (
	"go.uber.org/dig"
	"gopkg.in/urfave/cli.v1"
)

// Module of typical-go application
type Module struct {
	Name                string
	ShortName           string
	Usage               string
	ConfigPrefix        string
	Config              interface{}
	Command             *cli.Command
	LoadConfigFunc      interface{}
	OpenFunc            interface{}
	CloseFunc           interface{}
	SideEffects         []string
	AppSideEffects      []string
	TaskToolSideEffects []string
}

// Invoke the function for CLI command
func (m *Module) Invoke(invokeFunc interface{}) interface{} {
	return func(ctx *cli.Context) error {
		container := dig.New()
		container.Provide(m.LoadConfigFunc)
		container.Provide(m.OpenFunc)
		container.Provide(ctx.Args)
		return container.Invoke(invokeFunc)
	}
}
