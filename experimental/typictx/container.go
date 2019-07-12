package typictx

import "go.uber.org/dig"

// Container to return the depedency injection
func (c Context) Container() *dig.Container {
	container := dig.New()

	for _, contructor := range c.ArcheType.GetConstructors() {
		container.Provide(contructor)
	}

	for key := range c.Modules {
		module := c.Modules[key]
		container.Provide(module.LoadConfigFunc)
		container.Provide(module.OpenFunc)
	}

	return container
}
