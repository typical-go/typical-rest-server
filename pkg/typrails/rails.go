package typrails

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/utility/common"
	"github.com/typical-go/typical-go/pkg/utility/runn"
	"github.com/typical-go/typical-go/pkg/utility/runner"
	"github.com/typical-go/typical-rest-server/pkg/typrails/internal/tmpl"
	"github.com/urfave/cli/v2"
)

type rails struct {
	*typcore.Context
}

func (r *rails) scaffoldCmd() *cli.Command {
	return &cli.Command{
		Name:      "rails",
		Usage:     "Rails-like generator",
		ArgsUsage: "[table name]",
		Before: func(ctx *cli.Context) error {
			return common.LoadEnvFile()
		},
		Action: r.PreparedAction(r.scaffold),
	}
}

func (r *rails) scaffold(ctx *cli.Context, f Fetcher) (err error) {
	tableName := ctx.Args().First()
	if tableName == "" {
		return cli.ShowCommandHelp(ctx, "rails")
	}
	var e *Entity
	if e, err = f.Fetch(r.Package, tableName); err != nil {
		return
	}
	return generate(e)
}

func generate(e *Entity) error {
	repoPath := fmt.Sprintf("app/repository/%s_repo.go", e.Name)
	repoImplPath := fmt.Sprintf("app/repository/%s_repo_impl.go", e.Name)
	cachedRepoImplPath := fmt.Sprintf("app/repository/cached_%s_repo_impl.go", e.Name)
	servicePath := fmt.Sprintf("app/service/%s_service.go", e.Name)
	controllerPath := fmt.Sprintf("app/controller/%s_cntrl.go", e.Name)
	if common.IsFileExist(repoPath) {
		return fmt.Errorf("%s already exist", repoPath)
	}
	if common.IsFileExist(repoImplPath) {
		return fmt.Errorf("%s already exist", repoImplPath)
	}
	if common.IsFileExist(cachedRepoImplPath) {
		return fmt.Errorf("%s already exist", cachedRepoImplPath)
	}
	if common.IsFileExist(servicePath) {
		return fmt.Errorf("%s already exist", servicePath)
	}
	if common.IsFileExist(controllerPath) {
		return fmt.Errorf("%s already exist", controllerPath)
	}
	return runn.Execute(
		runner.NewWriteTemplate(repoPath, tmpl.Repo, e),
		runner.NewWriteTemplate(repoImplPath, tmpl.RepoImpl, e),
		runner.NewWriteTemplate(cachedRepoImplPath, tmpl.CachedRepoImpl, e),
		runner.NewWriteTemplate(servicePath, tmpl.Service, e),
		runner.NewWriteTemplate(controllerPath, tmpl.Controller, e),
		runner.NewGoFmt(repoPath, repoImplPath),
	)
}
