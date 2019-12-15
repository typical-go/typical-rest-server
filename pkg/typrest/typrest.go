package typrest

import (
	"fmt"
	"os"

	"github.com/typical-go/typical-go/pkg/typobj"
	"github.com/typical-go/typical-go/pkg/utility/runn"
	"github.com/typical-go/typical-go/pkg/utility/runner"
	"github.com/urfave/cli/v2"
)

// Module of typrest
type Module struct{}

// BuildCommands is commands to exectuce from Build-Tool
func (m *Module) BuildCommands(c typobj.Cli) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate CRUD (experimental)",
			Action:  m.scaffold,
		},
	}
}

func (m *Module) scaffold(ctx *cli.Context) (err error) {
	e := Entity{
		Name:           "music",
		TableName:      "musics",
		TypeName:       "Music",
		ProjectPackage: "github.com/typical-go/typical-rest-server",
	}
	repoPath := fmt.Sprintf("app/repository/%s_repo.go", e.Name)
	servicePath := fmt.Sprintf("app/service/%s_service.go", e.Name)
	os.Remove(repoPath)
	os.Remove(servicePath)
	return runn.Execute(
		runner.WriteTemplate{
			Target:   repoPath,
			Template: repoTmpl,
			Data:     e,
		},
		runner.WriteTemplate{
			Target:   servicePath,
			Template: serviceTmpl,
			Data:     e,
		},
	)
}
