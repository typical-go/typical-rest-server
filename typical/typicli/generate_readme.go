package typicli

import (
	"log"
	"os"
	"text/template"

	"gopkg.in/urfave/cli.v1"
)

func (t *TypicalCli) generateReadme(ctx *cli.Context) (err error) {

	readmeFile := t.ReadmeFileOrDefault()
	readmeTemplate := t.ReadmeTemplateOrDefault()

	templ, err := template.New("readme").Parse(readmeTemplate)

	if err != nil {
		return
	}

	file, err := os.Create(readmeFile)
	if err != nil {
		return
	}

	log.Printf("Generate ReadMe Document at '%s'", readmeFile)
	err = templ.Execute(file, Readme{
		Context: t.Context,
	})
	return nil
}
