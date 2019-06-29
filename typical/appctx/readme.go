package appctx

import (
	"bytes"
	"os"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/urfave/cli.v1"
)

// ConfigDoc for configuration documentation
func (c Context) ConfigDoc() string {
	buf := new(bytes.Buffer)
	for key := range c.Modules {
		module := c.Modules[key]
		buf.WriteString("\n")
		buf.WriteString(strcase.ToCamel(key))
		buf.WriteString("\n")

		envconfig.Usagef(module.ConfigPrefix(), module.Config(), buf, `
| Key | Type | Default | Request | Description |	
|---|---|---|---|---|	
{{range .}}|{{usage_key .}}|{{usage_type .}}|{{usage_default .}}|{{usage_required .}}|{{usage_description .}}|	
{{end}}`)
	}

	return buf.String()
}

func (c *Context) generateReadme(ctx *cli.Context) (err error) {
	t, err := template.New("readme").Parse(c.ReadmeTemplate)
	if err != nil {
		return
	}

	f, err := os.Create("README.md")
	if err != nil {
		return
	}

	err = t.Execute(f, c)
	return nil
}
