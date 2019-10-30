package typictx

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/slice"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
	"go.uber.org/dig"
)

// Context of typical application
type Context struct {
	typiobj.Modules
	Release
	Application  interface{}
	Name         string
	Description  string
	Root         string
	TestTargets  slice.Strings
	MockTargets  slice.Strings
	Constructors slice.Interfaces // TODO: remove this
	// *dig.Container
}

// Configurations return config list
func (c *Context) Configurations() (cfgs []typiobj.Configuration) {
	if configurer, ok := c.Application.(typiobj.Configurer); ok {
		cfgs = append(cfgs, configurer.Configure())
	}
	cfgs = append(cfgs, c.Modules.Configurations()...)
	return
}

// DockerCompose get docker compose
// func (c *Context) DockerCompose() (dockerCompose *docker.Compose) {
// 	dockerCompose = docker.NewCompose("3")
// 	for _, module := range c.Modules {
// 		moduleDocker := module.DockerCompose
// 		if moduleDocker == nil {
// 			continue
// 		}
// 		for _, name := range moduleDocker.ServiceKeys {
// 			dockerCompose.RegisterService(name, moduleDocker.Services[name])
// 		}
// 		for _, name := range moduleDocker.NetworkKeys {
// 			dockerCompose.RegisterNetwork(name, moduleDocker.Networks[name])
// 		}
// 		for _, name := range moduleDocker.VolumeKeys {
// 			dockerCompose.RegisterVolume(name, moduleDocker.Volumes[name])
// 		}
// 	}
// 	return
// }

// Destruct dependencies
func (c *Context) Destruct(container *dig.Container) (err error) {
	if destructor, ok := c.Application.(typiobj.Destructor); ok {
		if err = destructor.Destruct(container); err != nil {
			return
		}
	}
	return c.Modules.Destruct(container)
}

// Preparing context
// TODO: rename back to validate as conflicting with life cycle phase
func (c *Context) Preparing() (err error) {
	if err = c.validate(); err != nil {
		return
	}
	return
}

func (c *Context) validate() error {
	if c.Name == "" {
		return invalidContextError("Name can't not empty")
	}
	if c.Root == "" {
		return invalidContextError("Root can't not empty")
	}
	if _, ok := c.Application.(typiobj.Runner); !ok {
		return invalidContextError("Application must implement Runner")
	}
	return nil
}
