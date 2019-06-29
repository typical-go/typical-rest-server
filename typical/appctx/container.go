package appctx

import "go.uber.org/dig"

// Container to return the depedency injection
func (c Context) Container() *dig.Container {
	container := dig.New()

	for _, contructor := range c.Constructors {
		container.Provide(contructor)
	}

	for key := range c.Modules {
		module := c.Modules[key]
		for _, contructor := range module.Constructors() {
			container.Provide(contructor)
		}
	}

	return container
}
