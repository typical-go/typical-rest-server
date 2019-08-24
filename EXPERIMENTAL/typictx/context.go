package typictx

import (
	"os"
	"path/filepath"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/docker"
	"go.uber.org/dig"
)

// Context of typical application
type Context struct {
	Application
	Release

	Name        string
	Description string
	BinaryName  string

	Modules []*Module

	TestTargets []string
	MockTargets []string

	Constructors []interface{}

	container *dig.Container
}

// BinaryNameOrDefault return binary name of typiapp or default value
func (c *Context) BinaryNameOrDefault() string {
	if c.BinaryName == "" {
		dir, _ := os.Getwd()
		return filepath.Base(dir)
	}
	return c.BinaryName
}

// Invoke the function
func (c *Context) Invoke(function interface{}) error {
	if c.container == nil {
		c.container = dig.New()
		for _, module := range c.Modules {
			module.Inject(c.container)
		}
		for _, constructor := range c.Constructors {
			c.container.Provide(constructor)
		}
	}
	return c.container.Invoke(function)
}

// AddConstructor to add constructor
func (c *Context) AddConstructor(constructor interface{}) {
	c.Constructors = append(c.Constructors, constructor)
}

// AddMockTarget to add mock target
func (c *Context) AddMockTarget(mockTarget string) {
	c.MockTargets = append(c.MockTargets, mockTarget)
}

// AddTestTarget to add test target
func (c *Context) AddTestTarget(testTarget string) {
	c.TestTargets = append(c.TestTargets, testTarget)
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
