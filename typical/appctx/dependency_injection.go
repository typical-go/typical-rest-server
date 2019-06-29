package appctx

import (
	"go.uber.org/dig"
	"gopkg.in/urfave/cli.v1"
)

// DependencyInjection provide dependency injection handler
type DependencyInjection struct {
	constructors []interface{}
}

// NewDependencyInjection return new instance of DependencyInjection
func NewDependencyInjection(constructors ...interface{}) DependencyInjection {
	return DependencyInjection{
		constructors: constructors,
	}
}

// Constructors of dependency injectin
func (i *DependencyInjection) Constructors() []interface{} {
	return i.constructors
}

// Container of dependency injectin
func (i *DependencyInjection) Container() *dig.Container {
	container := dig.New()
	for _, contructor := range i.Constructors() {
		container.Provide(contructor)
	}
	return container
}

// InvokeFunction to invoke the function for CLI command
func (i *DependencyInjection) InvokeFunction(invokeFunc interface{}) interface{} {
	return func(ctx *cli.Context) error {
		container := i.Container()
		container.Provide(ctx.Args)
		return container.Invoke(invokeFunc)
	}
}
