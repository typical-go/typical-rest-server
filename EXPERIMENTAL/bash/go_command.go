package bash

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
)

// GoFmt for `go fmt`
func GoFmt(filename string) error {
	return exec.Command("go", "fmt", filename).Run()
}

// GoImports for `goimports`
func GoImports(filename string) error {
	return RunGoBin("goimports", "-w", filename)
}

// GoBuild for `go build`
func GoBuild(binaryName, mainPackage string, env ...string) error {
	args := []string{"build"}
	args = append(args, "-o", binaryName)
	args = append(args, "-ldflags", "-w -s")
	args = append(args, "./"+mainPackage)
	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

// GoTest for `go test` with coverprofile
func GoTest(targets []string) error {
	args := []string{"test"}
	args = append(args, targets...)
	args = append(args, "-coverprofile=cover.out")
	args = append(args, "-race")
	return Run("go", args...)
}

// GoGet for `go get`
func GoGet(packageName string) error {
	return Run("go", "get", packageName)
}

// RunGoBin to run binary in gobin folder
func RunGoBin(name string, args ...string) error {
	path := fmt.Sprintf("%s/%s/%s", build.Default.GOPATH, "bin", name)
	return Run(path, args...)
}

// GoClean for `go clean`
func GoClean(cleanArgs ...string) error {
	var args []string
	args = append(args, "clean")
	args = append(args, cleanArgs...)
	return Run("go", args...)
}
