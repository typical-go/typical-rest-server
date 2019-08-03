package typitask

import (
	"bytes"
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/readme"
)

const (
	configTemplate = `
| Key | Type | Default | Required | Description |	
|---|---|---|---|---|	
{{range .}}|{{usage_key .}}|{{usage_type .}}|{{usage_default .}}|{{usage_required .}}|{{usage_description .}}|	
{{end}}`
)

// GenerateReadme for generate typical applical readme
func GenerateReadme(ctx typictx.ActionContext) (err error) {
	recipe := readme.Recipe{
		Title:       ctx.Typical.Name,
		Description: ctx.Typical.Description,
		Sections: []readme.Section{
			{Title: "Getting Started", Content: gettingStartedInstruction()},
			{Title: "Usage", Content: usageInstruction()},
			{Title: "Build Tool", Content: buildToolInstruction()},
			{Title: "Make a release", Content: releaseInstruction()},
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

func gettingStartedInstruction() string {
	var md typirecipe.Markdown
	md.Writeln("This is intruction to start working with the project:")
	md.OrderedList(
		"Install [Go](https://golang.org/doc/install) or using homebrew if you're using macOS `brew install go`",
	)

	return md.String()
}

func usageInstruction() string {
	return `There is no specific requirement to run the application. `
}

func buildToolInstruction() string {
	return "Use `./typicalw` to execute development task"
}

func releaseInstruction() string {
	return "Use `./typicalw release -github=[TOKEN]` to make the release. You can found the release in `release` folder or github release page"
}
