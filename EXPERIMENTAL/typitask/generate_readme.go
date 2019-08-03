package typitask

import (
	"bytes"
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typigen/readme"
)

const (
	configTemplate = `
| Key | Type | Default | Required | Description |	
|---|---|---|---|---|	
{{range .}}|{{usage_key .}}|{{usage_type .}}|{{usage_default .}}|{{usage_required .}}|{{usage_description .}}|	
{{end}}`
	gettingStartedInstruction = `This is intruction to start working with the project:
1. Install go
2. Clone the project`
	usageInstruction   = `There is no specific requirement to run the application. `
	devToolInstruction = "Use `./typicalw` to execute development task"
	releaseInstruction = "Use `./typicalw release -github=[TOKEN]` to make the release.t r You can found the release in `release` folder or github release page"
)

// GenerateReadme for generate typical applical readme
func GenerateReadme(ctx typictx.ActionContext) (err error) {
	recipe := readme.ReadmeRecipe{
		Title:       ctx.Typical.Name,
		Description: ctx.Typical.Description,
		Sections: []readme.SectionPogo{
			{Title: "Getting Started", Content: gettingStartedInstruction},
			{Title: "Usage", Content: usageInstruction},
			{Title: "Development Tool", Content: devToolInstruction},
			{Title: "Make a release", Content: releaseInstruction},
			{Title: "Configurations", Content: configDoc(ctx.Typical)},
		},
	}
	log.Infof("Generate README.md")

	file, err := os.Create("README.md")
	if err != nil {
		return
	}
	defer file.Close()

	return recipe.Write(file)
}

func configDoc(ctx typictx.Context) string {
	buf := new(bytes.Buffer)

	for i, cfg := range ctx.Configs {
		if i > 0 {
			buf.WriteString("\n")
		}

		buf.WriteString(cfg.Description)
		buf.WriteString("\n")
		envconfig.Usagef(cfg.Prefix, cfg.Spec, buf, configTemplate)
	}

	return buf.String()
}
