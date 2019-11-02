package typictx

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/slice"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
	"github.com/urfave/cli"
)

// Context of typical application
type Context struct {
	Name         string
	Description  string
	Root         string
	Application  interface{}
	Modules      slice.Interfaces
	Release      Release
	TestTargets  slice.Strings
	MockTargets  slice.Strings
	Constructors slice.Interfaces
}

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

// BuildCommands return list of command for Build-Tool
func (c *Context) BuildCommands() (cmds []cli.Command) {
	if commandliner, ok := c.Application.(typiobj.BuildCLI); ok {
		cmds = append(cmds, commandliner.Command())
	}
	for _, module := range c.Modules {
		if commandliner, ok := module.(typiobj.BuildCLI); ok {
			cmds = append(cmds, commandliner.Command())
		}
	}
	return
}

// Provide the dependencies
func (c *Context) Provide() (constructors []interface{}) {
	constructors = append(constructors, c.Constructors...)
	if provider, ok := c.Application.(typiobj.Provider); ok {
		constructors = append(constructors, provider.Provide()...)
	}
	for _, module := range c.Modules {
		if provider, ok := module.(typiobj.Provider); ok {
			constructors = append(constructors, provider.Provide()...)
		}
	}
	return
}

// Destroy the dependencies
func (c *Context) Destroy() (destructors []interface{}) {
	if destroyer, ok := c.Application.(typiobj.Destroyer); ok {
		destructors = append(destructors, destroyer.Destroy()...)
	}
	for _, module := range c.Modules {
		if destroyer, ok := module.(typiobj.Destroyer); ok {
			destructors = append(destructors, destroyer.Destroy()...)
		}
	}
	return
}

// Prepare the run
func (c *Context) Prepare() (preparations []interface{}) {
	if preparer, ok := c.Application.(typiobj.Preparer); ok {
		preparations = append(preparations, preparer.Prepare()...)
	}
	return
}
