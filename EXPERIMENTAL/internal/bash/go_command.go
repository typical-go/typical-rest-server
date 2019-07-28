package bash

import (
	"fmt"
	"go/build"
	"strings"
)

var (
	goCommand = fmt.Sprintf("%s/bin/go", build.Default.GOROOT)
)

// GoFmt for `go fmt`
func GoFmt(filename string) error {
	return Silent(goCommand, "fmt", filename)
}

// GoImports for `goimports`
func GoImports(filename string) error {
	return RunGoBin("goimports", "-w", filename)
}

// GoBuild for `go build`
func GoBuild(binaryName, mainPackage string, ldflags ...string) error {
	args := []string{"build"}

	args = append(args, "-o", binaryName)
	args = append(args, "-ldflags", strings.Join(ldflags, " "))
	args = append(args, mainPackage)

	return Run(goCommand, args...)
}

// GoTest for `go test` with coverprofile
func GoTest(targets []string) error {
	args := []string{"test"}
	args = append(args, targets...)
	args = append(args, "-coverprofile=cover.out")
	return Run(goCommand, args...)
}

// GoGet for `go get`
func GoGet(packageName string) error {
	return Run(goCommand, "get", packageName)
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
	return Run(goCommand, args...)
}
