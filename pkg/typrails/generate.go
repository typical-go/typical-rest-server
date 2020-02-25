package typrails

import (
	"context"
	"fmt"
	"go/build"
	"os"
	"os/exec"

	"github.com/typical-go/typical-go/pkg/common"
	"github.com/typical-go/typical-go/pkg/runnerkit"
	"github.com/typical-go/typical-rest-server/pkg/typrails/internal/tmpl"
)

func generateController(ctx context.Context, e *Entity) (err error) {
	controllerPath := fmt.Sprintf("app/controller/%s_cntrl.go", e.Name)
	if common.IsFileExist(controllerPath) {
		return fmt.Errorf("%s already exist", controllerPath)
	}
	return runnerkit.Run(ctx,
		runnerkit.WriteTemplate(controllerPath, tmpl.Controller, e, 0666),
	)
}

func generateService(ctx context.Context, e *Entity) (err error) {
	servicePath := fmt.Sprintf("app/service/%s_service.go", e.Name)
	if common.IsFileExist(servicePath) {
		return fmt.Errorf("%s already exist", servicePath)
	}
	return runnerkit.Run(ctx,
		runnerkit.WriteTemplate(servicePath, tmpl.Service, e, 0666),
	)
}

func generateRepository(ctx context.Context, e *Entity) (err error) {
	repoPath := fmt.Sprintf("app/repository/%s.go", e.Name)
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
	return runnerkit.Run(ctx,
		runnerkit.WriteTemplate(repoPath, tmpl.Repo, e, 0666),
		runnerkit.WriteTemplate(repoImplPath, tmpl.RepoImpl, e, 0666),
		runnerkit.WriteTemplate(cachedRepoImplPath, tmpl.CachedRepoImpl, e, 0666),
		func() error {
			cmd := exec.Command(fmt.Sprintf("%s/bin/goimports", build.Default.GOPATH),
				"-w", repoPath, repoImplPath)
			cmd.Stderr = os.Stderr
			return cmd.Run()
		},
	)
}

func generateTransactional(ctx context.Context) (err error) {
	transactionalPath := "app/repository/transactional.go"
	transactionalTestPath := "app/repository/transactional_test.go"
	if common.IsFileExist(transactionalPath) {
		return nil
	}
	return runnerkit.Run(ctx,
		runnerkit.WriteString(transactionalPath, tmpl.Transactional, 0666),
		runnerkit.WriteString(transactionalTestPath, tmpl.TransactionalTest, 0666),
	)
}
