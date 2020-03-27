package typserver

import (
	"github.com/typical-go/typical-go/pkg/typapp"
)

// ServerModule handle dependency for server
type ServerModule struct{}

// Module return new instance of ServerModule
func Module() *ServerModule {
	return &ServerModule{}
}

// Provide dependency
func (*ServerModule) Provide() []*typapp.Constructor {
	return []*typapp.Constructor{
		typapp.NewConstructor(NewServer),
	}
}

// Destroy dependency
func (*ServerModule) Destroy() []*typapp.Destruction {
	return []*typapp.Destruction{
		typapp.NewDestruction(Shutdown),
	}
}
