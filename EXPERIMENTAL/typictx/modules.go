package typictx

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/slice"
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

// Modules is list of module
type Modules slice.Interfaces

// Configs return list of config
func (m Modules) Configs() (cfgs []Config) {
	for _, module := range m {
		if configurer, ok := module.(Configurer); ok {
			if cfg := configurer.Configure(); cfg != nil {
				cfgs = append(cfgs, cfg)
			}
		}
	}
	return
}

// Commands return list of command
func (m Modules) Commands() (cmds []cli.Command) {
	// TODO:
	return
}

// Construct dependency
func (m Modules) Construct(c *dig.Container) (err error) {
	for _, module := range m {
		if constructor, ok := module.(Constructor); ok {
			if err = constructor.Construct(c); err != nil {
				return
			}
		}
	}
	return
}

// Destruct dependency
func (m Modules) Destruct(c *dig.Container) (err error) {
	for _, module := range m {
		if constructor, ok := module.(Constructor); ok {
			if err = constructor.Destruct(c); err != nil {
				return
			}
		}
	}
	return
}

// Helps information
func (m Modules) Helps() (helps []Help) {
	for _, module := range m {
		var help Help
		help.Name = Name(module)
		help.Description = Description(module)
		// TODO: configuration
		helps = append(helps, help)
	}
	return
}
