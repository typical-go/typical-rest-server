package typictx

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/slice"
	"github.com/typical-go/typical-rest-server/pkg/utility/errkit"
	"go.uber.org/dig"
)

// Context of typical application
type Context struct {
	Application
	Modules
	Release
	Name         string
	Description  string
	Root         string
	TestTargets  slice.Strings
	MockTargets  slice.Strings
	Constructors slice.Interfaces // TODO: remove this
	*dig.Container
}

// Invoke the function
// TODO: remove this
func (c *Context) Invoke(function interface{}) (err error) {
	if c.Container == nil {
		c.Container = dig.New()
		c.Construct(c.Container)
	}
	return c.Container.Invoke(function)
}

// Configurations return config list
func (c *Context) Configurations() (cfgs []Configuration) {
	cfgs = append(cfgs, c.Application.Configure())
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

// Construct dependencies
// TODO: move to application
func (c *Context) Construct(container *dig.Container) (err error) {
	for _, constructor := range c.Constructors {
		if err = container.Provide(constructor); err != nil {
			return err
		}
	}
	return c.Modules.Construct(container)
}

// Destruct dependencies
// TODO: move to application
func (c *Context) Destruct(container *dig.Container) (err error) {
	var errs errkit.Errors
	if c.Application.StopFunc != nil {
		errs.Add(container.Invoke(c.Application.StopFunc))
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
	return nil
}
