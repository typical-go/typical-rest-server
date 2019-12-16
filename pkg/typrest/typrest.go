package typrest

import (
	"fmt"
	"os"

	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/utility/runn"
	"github.com/typical-go/typical-go/pkg/utility/runner"
	"github.com/urfave/cli/v2"
)

// Module of typrest
type Module struct{}

// BuildCommands is commands to exectuce from Build-Tool
func (m *Module) BuildCommands(c typcore.Cli) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate CRUD (experimental)",
			Action:  m.scaffold,
		},
	}
}

// ID        int64     `json:"id"`
// Title     string    `json:"title" validate:"required"`
// Author    string    `json:"author" validate:"required"`
// UpdatedAt time.Time `json:"updated_at"`
// CreatedAt time.Time `json:"created_at"`

func (m *Module) scaffold(ctx *cli.Context) (err error) {
	e := Entity{
		Name:           "music",
		Table:          "musics",
		Type:           "Music",
		Cache:          "MUSIC",
		ProjectPackage: "github.com/typical-go/typical-rest-server",
		Fields: []Field{
			{Name: "ID", Type: "int64"},
			{Name: "Title", Type: "string"},
			{Name: "UpdatedAt", Type: "time.Time"},
			{Name: "CreatedAt", Type: "time.Time"},
		},
	}
	repoPath := fmt.Sprintf("app/repository/%s_repo.go", e.Name)
	repoImplPath := fmt.Sprintf("app/repository/%s_repo_impl.go", e.Name)
	cachedRepoImplPath := fmt.Sprintf("app/repository/cached_%s_repo_impl.go", e.Name)
	servicePath := fmt.Sprintf("app/service/%s_service.go", e.Name)
	controllerPath := fmt.Sprintf("app/controller/%s_cntrl.go", e.Name)
	os.Remove(repoPath)
	os.Remove(repoImplPath)
	os.Remove(cachedRepoImplPath)
	os.Remove(servicePath)
	os.Remove(controllerPath)
	return runn.Execute(
		runner.WriteTemplate{Target: repoPath, Template: repoTemplate, Data: e},
		runner.WriteTemplate{Target: repoImplPath, Template: repoImplTemplate, Data: e},
		runner.WriteTemplate{Target: cachedRepoImplPath, Template: cachedRepoImplTemplate, Data: e},
		runner.WriteTemplate{Target: servicePath, Template: serviceTemplate, Data: e},
		runner.WriteTemplate{Target: controllerPath, Template: constrollerTemplate, Data: e},
	)
}
