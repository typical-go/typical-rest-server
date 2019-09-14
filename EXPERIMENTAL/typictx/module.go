package typictx

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/docker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/slice"
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

// Module of typical-go application
type Module struct {
	Config
	Name          string
	OpenFunc      interface{}
	CloseFunc     interface{}
	Command       *Command
	DockerCompose *docker.Compose
	Constructors  slice.Interfaces
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
// TODO: revisit this as separate container is wrong implementation
func (m *Module) Invoke(invokeFunc interface{}) interface{} {
	return func(ctx *cli.Context) error {
		container := dig.New()
		container.Provide(ctx.Args) // NOTE: inject cli arguments
		m.Inject(container)
		return container.Invoke(invokeFunc)
	}
}

// GetName to get name
func (m *Module) GetName() string {
	return m.Name
}
