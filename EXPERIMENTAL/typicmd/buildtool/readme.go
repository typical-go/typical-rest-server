package buildtool

import (
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/buildtool/readme"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

const (
	readmeFile     = "README.md"
	configTemplate = `| Key | Type | Default | Required | Description |	
|---|---|---|---|---|{{range .}}
|{{usage_key .}}|{{usage_type .}}|{{usage_default .}}|{{usage_required .}}|{{usage_description .}}|{{end}}`
)

// GenerateReadme for generate typical applical readme
func generateReadme(a *typictx.ActionContext) (err error) {
	readme0 := readme.DefaultReadme().
		SetTitle(a.Name).
		SetDescription(a.Description).
		SetSection("Configuration", func(md *readme.Markdown) (err error) {
			for _, module := range a.Context.Modules {
				name := strcase.ToCamel(module.Name)
				if name != "" {
					md.Heading3(name)
				}
				var builder strings.Builder
				envconfig.Usagef(module.Prefix, module.Spec, &builder, configTemplate)
				md.Writeln(builder.String())
			}
			return
		})
	log.Infof("Generate new %s", readmeFile)
	err = readme0.OutputToFile(readmeFile)
	if err != nil {
		return
	}
	return
}
