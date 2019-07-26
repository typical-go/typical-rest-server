package typitask

import (
	"bytes"
	"os"
	"text/template"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

const configTemplate = `
| Key | Type | Default | Required | Description |	
|---|---|---|---|---|	
{{range .}}|{{usage_key .}}|{{usage_type .}}|{{usage_default .}}|{{usage_required .}}|{{usage_description .}}|	
{{end}}`

// GenerateReadme for generate typical applical readme
func GenerateReadme(ctx typictx.ActionContext) (err error) {
	readmeFile := ctx.Typical.ReadmeFileOrDefault()
	readmeTemplate := ctx.Typical.ReadmeTemplateOrDefault()

	templ, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		return
	}

	file, err := os.Create(readmeFile)
	if err != nil {
		return
	}

	log.Infof("Generate ReadMe Document at '%s'", readmeFile)
	return templ.Execute(file, Readme{
		Context: ctx.Typical,
	})
}

// Readme represent readme structured data
type Readme struct {
	typictx.Context
}

// ConfigDoc for configuration documentation
func (r Readme) ConfigDoc() string {
	buf := new(bytes.Buffer)

	for _, cfg := range r.Context.Configs {
		buf.WriteString("\n")
		buf.WriteString(cfg.Description)
		buf.WriteString("\n")
		envconfig.Usagef(cfg.Prefix, cfg.Spec, buf, configTemplate)
	}

	return buf.String()
}
