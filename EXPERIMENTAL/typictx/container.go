package typictx

import "go.uber.org/dig"

// Container to return the depedency injection
func (c Context) Container() *dig.Container {
	container := dig.New()

	for _, constructor := range c.ArcheType.GetConstructors() {
		container.Provide(constructor)
	}

	for _, module := range c.Modules {
		module.Inject(container)
	}

	return container
}
