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

func generate(e *Entity) (err error) {
	path := "app/repository"
	repoPath := fmt.Sprintf("%s/%s_repo.go", path, e.Name)
	repoImplPath := fmt.Sprintf("%s/%s_repo_impl.go", path, e.Name)
	cachedRepoImplPath := fmt.Sprintf("%s/cached_%s_repo_impl.go", path, e.Name)
	servicePath := fmt.Sprintf("%s/%s_service.go", path, e.Name)
	controllerPath := fmt.Sprintf("%s/%s_cntrl.go", path, e.Name)
	transactionalPath := fmt.Sprintf("%s/transactional.go", path)
	transactionalTestPath := fmt.Sprintf("%s/transactional_test.go", path)
	if !common.IsFileExist(transactionalPath) {
		if err = runn.Execute(
			runner.NewWriteString(transactionalPath, tmpl.Transactional),
			runner.NewWriteString(transactionalTestPath, tmpl.TransactionalTest),
		); err != nil {
			return err
		}
	}
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
