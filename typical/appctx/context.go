package appctx

import (
	"bytes"

	"go.uber.org/dig"

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
