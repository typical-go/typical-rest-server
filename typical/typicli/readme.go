package typicli

import (
	"bytes"

	"github.com/iancoleman/strcase"
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/typical/appctx"
)

const configTemplate = `
| Key | Type | Default | Request | Description |	
|---|---|---|---|---|	
{{range .}}|{{usage_key .}}|{{usage_type .}}|{{usage_default .}}|{{usage_required .}}|{{usage_description .}}|	
{{end}}`

// Readme represent readme structured data
type Readme struct {
	appctx.Context
}

// ConfigDoc for configuration documentation
func (r Readme) ConfigDoc() string {
	buf := new(bytes.Buffer)

	buf.WriteString("\nApplication\n")
	envconfig.Usagef(r.TypiApp.ConfigPrefix, r.TypiApp.Config, buf, configTemplate)

	for i := range r.Modules {
		module := r.Modules[i]
		buf.WriteString("\n")
		buf.WriteString(strcase.ToCamel(module.Name))
		buf.WriteString("\n")

		envconfig.Usagef(module.ConfigPrefix, module.Config, buf, configTemplate)
	}

	return buf.String()
}
