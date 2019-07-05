package typicli

import (
	"os"
	"text/template"

	"gopkg.in/urfave/cli.v1"
)

func (t *TypicalCli) generateReadme(ctx *cli.Context) (err error) {
	templ, err := template.New("readme").Parse(t.ReadmeTemplate)
	if err != nil {
		return
	}

	file, err := os.Create("README.md")
	if err != nil {
		return
	}

	err = templ.Execute(file, Readme{
		Context: t.Context,
	})
	return nil
}
