package appctx

import (
	"go.uber.org/dig"
	"gopkg.in/urfave/cli.v1"
)

// // Module in typical-go applicaiton
// type Module interface {
// 	Name() string
// 	ShortName() string
// 	ConfigPrefix() string
// 	Config() interface{}
// 	Constructors() []interface{}
// 	Command() cli.Command
// 	LoadFunc() interface{}
// }

// Module of typical-go application
type Module struct {
	Name           string
	ShortName      string
	Usage          string
	ConfigPrefix   string
	Config         interface{}
	Command        *cli.Command
	LoadConfigFunc interface{}
	OpenFunc       interface{}
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
