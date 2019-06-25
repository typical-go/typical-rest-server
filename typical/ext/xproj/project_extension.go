package xproj

import (
	"fmt"
	"os"
	"text/template"

	"github.com/typical-go/typical-rest-server/typical/appctx"
	"github.com/typical-go/typical-rest-server/typical/ext"
	"gopkg.in/urfave/cli.v1"
)

// ProjectExtension provide standard command to see project context and configuration
type ProjectExtension struct {
	ext.Extension
	ext.ActionTrigger
}

// NewProjectExtension return new instance of ProjectExtension
func NewProjectExtension(context appctx.Context) *ProjectExtension {
	return &ProjectExtension{
		ActionTrigger: ext.ActionTrigger{Context: context},
	}
}

// Setup go extension
func (e *ProjectExtension) Setup() error {
	return fmt.Errorf("not implement")
}

//Command for go extension
func (e *ProjectExtension) Command() cli.Command {
	return cli.Command{
		Name:      "project",
		ShortName: "proj",
		Subcommands: []cli.Command{
			{Name: "config", Usage: "Config details", Action: e.Print(configDetail)},
			{Name: "context", Usage: "Context details", Action: e.Print(contextDetail)},
			{Name: "readme", Usage: "Generate readme", Action: e.generateReadme},
		},
	}
}

func (e *ProjectExtension) generateReadme(ctx *cli.Context) (err error) {
	t, err := template.New("readme").Parse(e.Context.ReadmeTemplate)
	if err != nil {
		return
	}

	f, err := os.Create("README.md")
	if err != nil {
		return
	}

	err = t.Execute(f, e.Context)
	return nil
}
