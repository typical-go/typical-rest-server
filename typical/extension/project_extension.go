package extension

import (
	"fmt"

	"github.com/typical-go/typical-rest-server/typical/task/generate"
	"github.com/typical-go/typical-rest-server/typical/task/project"
	"gopkg.in/urfave/cli.v1"
)

// ProjectExtension provide standard command to see project context and configuration
type ProjectExtension struct{}

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
			{Name: "config", Usage: "Config details", Action: print(project.ConfigDetail)},
			{Name: "context", Usage: "Context details", Action: print(project.ContextDetail)},
			{Name: "readme", Usage: "Generate readme", Action: run(generate.Readme)},
		},
	}
}
