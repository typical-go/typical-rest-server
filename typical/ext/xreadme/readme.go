package xreadme

import (
	"bytes"

	"github.com/iancoleman/strcase"
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/typical/appctx"
)

// Readme represent readme structured data
type Readme struct {
	appctx.Context
}

// ConfigDoc for configuration documentation
func (r Readme) ConfigDoc() string {
	buf := new(bytes.Buffer)
	for key := range r.Modules {
		module := r.Modules[key]
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
