package typictx

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/slice"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
)

// Context of typical application
type Context struct {
	Name         string
	Description  string
	Root         string
	Application  interface{}
	Modules      typiobj.Modules
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
