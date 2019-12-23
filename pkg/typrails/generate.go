package typrails

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"

	"github.com/typical-go/typical-go/pkg/utility/common"
	"github.com/typical-go/typical-go/pkg/utility/runn"
	"github.com/typical-go/typical-go/pkg/utility/runner"
	"github.com/typical-go/typical-rest-server/pkg/typrails/internal/tmpl"
)

func generateController(e *Entity) (err error) {
	controllerPath := fmt.Sprintf("app/controller/%s_cntrl.go", e.Name)
	if common.IsFileExist(controllerPath) {
		return fmt.Errorf("%s already exist", controllerPath)
	}
	return runn.Execute(
		runner.NewWriteTemplate(controllerPath, tmpl.Controller, e),
	)
}

func generateService(e *Entity) (err error) {
	servicePath := fmt.Sprintf("app/service/%s_service.go", e.Name)
	if common.IsFileExist(servicePath) {
		return fmt.Errorf("%s already exist", servicePath)
	}
	return runn.Execute(
		runner.NewWriteTemplate(servicePath, tmpl.Service, e),
	)
}

func generateRepository(e *Entity) (err error) {
	repoPath := fmt.Sprintf("app/repository/%s_repo.go", e.Name)
	repoImplPath := fmt.Sprintf("app/repository/%s_repo_impl.go", e.Name)
	cachedRepoImplPath := fmt.Sprintf("app/repository/cached_%s_repo_impl.go", e.Name)
	if common.IsFileExist(repoPath) {
		return fmt.Errorf("%s already exist", repoPath)
	}
	if common.IsFileExist(repoImplPath) {
		return fmt.Errorf("%s already exist", repoImplPath)
	}
	if common.IsFileExist(cachedRepoImplPath) {
		return fmt.Errorf("%s already exist", cachedRepoImplPath)
	}
	return runn.Execute(
		runner.NewWriteTemplate(repoPath, tmpl.Repo, e),
		runner.NewWriteTemplate(repoImplPath, tmpl.RepoImpl, e),
		runner.NewWriteTemplate(cachedRepoImplPath, tmpl.CachedRepoImpl, e),
		func() error {
			cmd := exec.Command(fmt.Sprintf("%s/bin/goimports", build.Default.GOPATH),
				"-w", repoPath, repoImplPath)
			cmd.Stderr = os.Stderr
			return cmd.Run()
		},
	)
}

func generateTransactional() (err error) {
	transactionalPath := "app/repository/transactional.go"
	transactionalTestPath := "app/repository/transactional_test.go"
	if common.IsFileExist(transactionalPath) {
		return nil
	}
	return runn.Execute(
		runner.NewWriteString(transactionalPath, tmpl.Transactional),
		runner.NewWriteString(transactionalTestPath, tmpl.TransactionalTest),
	)
}
