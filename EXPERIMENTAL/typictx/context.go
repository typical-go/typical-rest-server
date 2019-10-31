package typictx

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/slice"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
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

// Validate context
func (c *Context) Validate() error {
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
