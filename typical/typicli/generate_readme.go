package typicli

import (
	"os"
	"text/template"

	"github.com/typical-go/typical-rest-server/typical/ext/xreadme"
	"gopkg.in/urfave/cli.v1"
)

func (t *Typical) generateReadme(ctx *cli.Context) (err error) {
	templ, err := template.New("readme").Parse(t.ReadmeTemplate)
	if err != nil {
		return
	}

	file, err := os.Create("README.md")
	if err != nil {
		return
	}

	err = templ.Execute(file, xreadme.Readme{
		Context: t.Context,
	})
	return nil
}
