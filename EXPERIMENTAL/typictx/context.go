package typictx

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/docker"

	log "github.com/sirupsen/logrus"

	"go.uber.org/dig"
)

// Context of typical application
type Context struct {
	Application
	Release

	Name        string
	Version     string
	Description string
	BinaryName  string

	Modules []*Module

	TestTargets []string
	MockTargets []string

	Constructors []interface{}

	DockerCompose *docker.Compose
	container     *dig.Container
}

// BinaryNameOrDefault return binary name of typiapp or default value
func (c *Context) BinaryNameOrDefault() string {
	if c.BinaryName == "" {
		dir, _ := os.Getwd()
		chunks := strings.Split(dir, "/")
		return chunks[len(chunks)-1]
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

// InvokeInitiation to invoke initiation functions
func (c *Context) InvokeInitiation() (err error) {
	for _, initiation := range c.Initiations {
		log.Info("Invoke initiation: " + reflect.TypeOf(initiation).String())
		err = c.Invoke(initiation)
		if err != nil {
			return
		}
	}
	return
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

// ModulesWithConfig return module that have config
func (c *Context) ModulesWithConfig() (modules []*Module) {
	for _, module := range c.Modules {
		if module.ConfigSpec != nil {
			modules = append(modules, module)
		}
	}
	return
}

// ReleaseVersion to return release version
func (c *Context) ReleaseVersion() (version string) {
	version = fmt.Sprintf("v%s", c.Version)
	if c.Release.Alpha {
		version = fmt.Sprintf("%s-alpha", version)
	}
	return
}

// Deadline implementation
func (*Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done implementation
func (*Context) Done() <-chan struct{} {
	return nil
}

// Err implementation
func (*Context) Err() error {
	return nil
}

// Value implementation
func (*Context) Value(key interface{}) interface{} {
	return nil
}
