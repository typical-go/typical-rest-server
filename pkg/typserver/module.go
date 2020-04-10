package typserver

import (
	"github.com/typical-go/typical-go/pkg/typapp"
)

var (
	_ typapp.Provider  = (*ServerModule)(nil)
	_ typapp.Destroyer = (*ServerModule)(nil)
)

// ServerModule handle dependency for server
type ServerModule struct{}

// Module return new instance of ServerModule
func Module() *ServerModule {
	return &ServerModule{}
}

// Constructors of module
func (*ServerModule) Constructors() []*typapp.Constructor {
	return []*typapp.Constructor{
		typapp.NewConstructor(NewServer),
	}
}

// Destructions of module
func (*ServerModule) Destructions() []*typapp.Destruction {
	return []*typapp.Destruction{
		typapp.NewDestruction(Shutdown),
	}
}
