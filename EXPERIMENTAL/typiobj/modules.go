package typiobj

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/slice"
	"github.com/urfave/cli"
)

// Modules is list of module
type Modules slice.Interfaces

// Commands return list of command
func (m Modules) Commands() (cmds []cli.Command) {
	for _, module := range m {
		if cmdline, ok := module.(CommandLiner); ok {
			cmds = append(cmds, cmdline.CommandLine())
		}
	}
	return
}

// Provide dependency
func (m Modules) Provide() (constructors []interface{}) {
	for _, module := range m {
		if provider, ok := module.(Provider); ok {
			constructors = append(constructors, provider.Provide()...)
		}
	}
	return
}

// Destroy dependency
func (m Modules) Destroy() (destructors []interface{}) {
	for _, module := range m {
		if destroyer, ok := module.(Destroyer); ok {
			destructors = append(destructors, destroyer.Destroy()...)
		}
	}
	return
}
