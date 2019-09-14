package typictx

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/docker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/slice"
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

// GetName to get name
func (m *Module) GetName() string {
	return m.Name
}
