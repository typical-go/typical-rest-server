package appctx

import "go.uber.org/dig"

// Container to return the depedency injection
func (c Context) Container() *dig.Container {
	container := dig.New()

	container.Provide(c.ConfigLoader.LoadFunc())

	for _, contructor := range c.Constructors {
		container.Provide(contructor)
	}

	for key := range c.Modules {
		module := c.Modules[key]
		container.Provide(module.LoadFunc())
		for _, contructor := range module.Constructors() {
			container.Provide(contructor)
		}
	}

	return container
}
