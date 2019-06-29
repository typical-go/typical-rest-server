package appctx

import (
	"bytes"

	"go.uber.org/dig"
	"gopkg.in/urfave/cli.v1"

	"github.com/iancoleman/strcase"
	"github.com/kelseyhightower/envconfig"
)

// Context of typical application
type Context struct {
	Name         string
	ConfigPrefix string
	Path         string
	Version      string
	Description  string
	// Container      *dig.Container
	Constructors   []interface{}
	Modules        map[string]Module
	ReadmeTemplate string
}

// Container to return the depedency injection
func (c Context) Container() *dig.Container {
	container := dig.New()

	for _, contructor := range c.Constructors {
		container.Provide(contructor)
	}

	return container
}

// ConfigDoc for configuration documentation
func (c Context) ConfigDoc() string {
	buf := new(bytes.Buffer)
	for key := range c.Modules {
		buf.WriteString("\n")
		buf.WriteString(strcase.ToCamel(key))
		buf.WriteString("\n")
		envconfig.Usagef(c.ConfigPrefix, c.Modules[key].Config, buf, `
| Key | Type | Default | Request | Description |	
|---|---|---|---|---|	
{{range .}}|{{usage_key .}}|{{usage_type .}}|{{usage_default .}}|{{usage_required .}}|{{usage_description .}}|	
{{end}}`)
	}

	return buf.String()
}

// Run to start the command line interface
func (c *Context) Run(arguments []string) error {
	app := cli.NewApp()
	app.Name = c.Name
	app.Usage = ""
	app.Description = c.Description
	app.Version = c.Version

	for key := range c.Modules {
		app.Commands = append(app.Commands, c.Modules[key].Command())
	}
	return app.Run(arguments)
}
