package appctx

import "go.uber.org/dig"

// Context of typical application
type Context struct {
	TypiApp
	TypiCli

	Name        string
	Version     string
	Description string

	ReadmeTemplate string
	ReadmeFile     string

	Modules []*Module
}

// ReadmeTemplateOrDefault return readme template field or the default template
func (c Context) ReadmeTemplateOrDefault() string {
	if c.ReadmeTemplate == "" {
		return defaultReadmeTemplate
	}
	return c.ReadmeTemplate
}

// ReadmeFileOrDefault return readme file field or default template
func (c Context) ReadmeFileOrDefault() string {
	if c.ReadmeFile == "" {
		return defaultReadmeFile
	}

	return c.ReadmeFile
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
