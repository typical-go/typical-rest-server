package xproj

import (
	"fmt"

	"github.com/typical-go/typical-rest-server/typical/ext"
	"gopkg.in/urfave/cli.v1"
)

// ProjectExtension provide standard command to see project context and configuration
type ProjectExtension struct {
	ext.Extension
	ext.ActionTrigger
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
			{Name: "readme", Usage: "Generate readme", Action: e.Run(generateReadme)},
		},
	}
}
