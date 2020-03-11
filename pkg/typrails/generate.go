package typrails

import (
	"context"
	"fmt"
	"go/build"
	"os"
	"os/exec"

	"github.com/typical-go/typical-go/pkg/exor"
	"github.com/typical-go/typical-rest-server/pkg/typrails/internal/tmpl"
)

func generateController(ctx context.Context, e *Entity) (err error) {
	controllerPath := fmt.Sprintf("app/controller/%s_cntrl.go", e.Name)
	if isFileExist(controllerPath) {
		return fmt.Errorf("%s already exist", controllerPath)
	}
	return exor.Execute(ctx,
		exor.NewWriteTemplate(controllerPath, tmpl.Controller, e),
	)
}

func generateService(ctx context.Context, e *Entity) (err error) {
	servicePath := fmt.Sprintf("app/service/%s_service.go", e.Name)
	if isFileExist(servicePath) {
		return fmt.Errorf("%s already exist", servicePath)
	}
	return exor.Execute(ctx,
		exor.NewWriteTemplate(servicePath, tmpl.Service, e),
	)
}

func generateRepository(ctx context.Context, e *Entity) (err error) {
	repoPath := fmt.Sprintf("app/repository/%s.go", e.Name)
	repoImplPath := fmt.Sprintf("app/repository/%s_repo_impl.go", e.Name)
	cachedRepoImplPath := fmt.Sprintf("app/repository/cached_%s_repo_impl.go", e.Name)
	if isFileExist(repoPath) {
		return fmt.Errorf("%s already exist", repoPath)
	}
	if isFileExist(repoImplPath) {
		return fmt.Errorf("%s already exist", repoImplPath)
	}
	if isFileExist(cachedRepoImplPath) {
		return fmt.Errorf("%s already exist", cachedRepoImplPath)
	}
	return exor.Execute(ctx,
		exor.NewWriteTemplate(repoPath, tmpl.Repo, e),
		exor.NewWriteTemplate(repoImplPath, tmpl.RepoImpl, e),
		exor.NewWriteTemplate(cachedRepoImplPath, tmpl.CachedRepoImpl, e),
		exor.New(func(ctx context.Context) error {
			cmd := exec.CommandContext(ctx, fmt.Sprintf("%s/bin/goimports", build.Default.GOPATH),
				"-w", repoPath, repoImplPath)
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}),
	)
}

func generateTransactional(ctx context.Context) (err error) {
	transactionalPath := "app/repository/transactional.go"
	transactionalTestPath := "app/repository/transactional_test.go"
	if isFileExist(transactionalPath) {
		return nil
	}
	return exor.Execute(ctx,
		exor.NewWriteString(transactionalPath, tmpl.Transactional),
		exor.NewWriteString(transactionalTestPath, tmpl.TransactionalTest),
	)
}

func isFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
