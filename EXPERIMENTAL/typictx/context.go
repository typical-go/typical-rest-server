package typictx

import (
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/docker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/slice"
	"go.uber.org/dig"
)

// Context of typical application
type Context struct {
	Application
	Release
	Name         string
	Description  string
	Root         string
	Modules      []*Module
	TestTargets  slice.Strings
	MockTargets  slice.Strings
	Constructors slice.Interfaces
	container    *dig.Container
}

// Invoke the function
func (c *Context) Invoke(function interface{}) error {
	if c.container == nil {
		c.container = dig.New()
		for _, module := range c.Modules {
			module.Inject(c.container)
		}
		for _, constructor := range c.Constructors {
			err := c.container.Provide(constructor)
			if err != nil {
				return err
			}
		}
	}
	return c.container.Invoke(function)
}

// AddConstructor to add constructor
func (c *Context) AddConstructor(constructor interface{}) {
	c.Constructors = append(c.Constructors, constructor)
}

// ConfigAccessors return list of config accessor
func (c *Context) ConfigAccessors() (accessors []ConfigAccessor) {
	accessors = append(accessors, &c.Application)
	for _, module := range c.Modules {
		if module.Spec != nil {
			accessors = append(accessors, module)
		}
	}
	return
}

// DockerCompose get docker compose
func (c *Context) DockerCompose() (dockerCompose *docker.Compose) {
	dockerCompose = docker.NewCompose("3")
	for _, module := range c.Modules {
		moduleDocker := module.DockerCompose
		if moduleDocker == nil {
			continue
		}
		for _, name := range moduleDocker.ServiceKeys {
			dockerCompose.RegisterService(name, moduleDocker.Services[name])
		}
		for _, name := range moduleDocker.NetworkKeys {
			dockerCompose.RegisterNetwork(name, moduleDocker.Networks[name])
		}
		for _, name := range moduleDocker.VolumeKeys {
			dockerCompose.RegisterVolume(name, moduleDocker.Volumes[name])
		}
	}
	return
}

// Validate the context
func (c *Context) Validate() error {
	log.Debug("Validate the context")
	if c.Name == "" {
		return invalidContextError("Name can't not empty")
	}
	if c.Root == "" {
		return invalidContextError("Root can't not empty")
	}
	return nil
}
