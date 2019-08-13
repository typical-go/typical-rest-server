package typictx

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

// Module of typical-go application
type Module struct {
	Name string

	ConfigPrefix string
	ConfigSpec   interface{}

	OpenFunc  interface{}
	CloseFunc interface{}

	Command       *Command
	DockerCompose *DockerCompose

	Constructors []interface{}
	SideEffects  []*SideEffect
}

// Inject dependencies for the module
func (m *Module) Inject(container *dig.Container) {
	for _, constructor := range m.Constructors {
		container.Provide(constructor)
	}
	container.Provide(m.OpenFunc)
	return
}

// Invoke the function for CLI command
func (m *Module) Invoke(invokeFunc interface{}) interface{} {
	return func(ctx *cli.Context) error {
		container := dig.New()
		container.Provide(ctx.Args) // NOTE: inject cli arguments
		m.Inject(container)

		return container.Invoke(invokeFunc)
	}
}

// CamelConfigPrefix return config prefix in camel case
func (m *Module) CamelConfigPrefix() string {
	if m.ConfigPrefix == "" {
		return ""
	}
	return strcase.ToCamel(strings.ToLower(m.ConfigPrefix))
}
