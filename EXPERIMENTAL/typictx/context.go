package typictx

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/collection"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
	"github.com/urfave/cli"
)

// Context of typical application
type Context struct {
	Name         string
	Description  string
	Root         string
	AppModule    interface{}
	Modules      collection.Interfaces
	Release      Release
	TestTargets  collection.Strings
	MockTargets  collection.Strings
	Constructors collection.Interfaces
}

// Validate context
func (c *Context) Validate() error {
	if c.Name == "" {
		return invalidContextError("Name can't not empty")
	}
	if c.Root == "" {
		return invalidContextError("Root can't not empty")
	}
	if _, ok := c.AppModule.(typiobj.Runner); !ok {
		return invalidContextError("Application must implement Runner")
	}
	return nil
}

// AllModule return app module and modules
func (c *Context) AllModule() (modules []interface{}) {
	modules = append(modules, c.AppModule)
	modules = append(modules, c.Modules...)
	return
}

// BuildCommands return list of command for Build-Tool
func (c *Context) BuildCommands() (cmds []cli.Command) {
	for _, module := range c.AllModule() {
		if commander, ok := module.(BuildCommander); ok {
			cmds = append(cmds, commander.BuildCommand(c))
		}
	}
	return
}

// Provide the dependencies
func (c *Context) Provide() (constructors []interface{}) {
	constructors = append(constructors, c.Constructors...)
	for _, module := range c.AllModule() {
		if provider, ok := module.(typiobj.Provider); ok {
			constructors = append(constructors, provider.Provide()...)
		}
	}
	return
}

// Destroy the dependencies
func (c *Context) Destroy() (destructors []interface{}) {
	for _, module := range c.AllModule() {
		if destroyer, ok := module.(typiobj.Destroyer); ok {
			destructors = append(destructors, destroyer.Destroy()...)
		}
	}
	return
}

// Prepare the run
func (c *Context) Prepare() (preparations []interface{}) {
	for _, module := range c.AllModule() {
		if preparer, ok := module.(typiobj.Preparer); ok {
			preparations = append(preparations, preparer.Prepare()...)
		}
	}
	return
}
