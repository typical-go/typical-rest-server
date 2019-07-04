package appctx

import "go.uber.org/dig"

// Context of typical application
type Context struct {
	TypiApp
	TypiCli

	Name        string
	Version     string
	Description string

	Modules        []*Module
	ReadmeTemplate string
}

// Container to return the depedency injection
func (c Context) Container() *dig.Container {
	container := dig.New()

	container.Provide(c.TypiApp.ConfigLoadFunc)

	for _, contructor := range c.TypiApp.Constructors {
		container.Provide(contructor)
	}

	for key := range c.Modules {
		module := c.Modules[key]
		container.Provide(module.LoadConfigFunc)
		container.Provide(module.OpenFunc)
	}

	return container
}
