package typictx

import (
	"os"
	"strings"

	"go.uber.org/dig"
)

// Context of typical application
type Context struct {
	Name        string
	Version     string
	Description string

	BinaryName string

	App Application

	ReadmeTemplate string
	ReadmeFile     string

	Configs []Config
	Modules []*Module

	Constructors []interface{}
	TestTargets  []string
	MockTargets  []string

	Github *Github
}

// ReadmeTemplateOrDefault return readme template field or the default template
func (c *Context) ReadmeTemplateOrDefault() string {
	if c.ReadmeTemplate == "" {
		return defaultReadmeTemplate
	}
	return c.ReadmeTemplate
}

// ReadmeFileOrDefault return readme file field or default template
func (c *Context) ReadmeFileOrDefault() string {
	if c.ReadmeFile == "" {
		return defaultReadmeFile
	}

	return c.ReadmeFile
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

// Container to return the depedency injection
func (c *Context) Container() *dig.Container {
	container := dig.New()

	for _, constructor := range c.Constructors {
		container.Provide(constructor)
	}

	for _, constructor := range c.Constructors {
		container.Provide(constructor)
	}

	for _, module := range c.Modules {
		module.Inject(container)
	}

	return container
}

// CheckModuleStatus to check module availability status
func (c *Context) CheckModuleStatus() (statusReport map[string]string) {
	statusReport = make(map[string]string)
	for _, module := range c.Modules {
		err := c.Container().Invoke(module.StatusFunc)
		if err != nil {
			statusReport[module.Name] = err.Error()
		} else {
			statusReport[module.Name] = "ok"
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
