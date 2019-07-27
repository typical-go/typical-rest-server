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
func GoFmt(filename string) {
	Silent(goCommand, "fmt", filename)
}

// GoImports for `goimports`
func GoImports(filename string) {
	RunGoBin("goimports", "-w", filename)
}

// GoBuild for `go build`
func GoBuild(binaryName, mainPackage string, ldflags ...string) {
	args := []string{"build"}

	args = append(args, "-o", binaryName)
	args = append(args, "-ldflags", strings.Join(ldflags, " "))
	args = append(args, mainPackage)

	Run(goCommand, args...)
}

// GoTest for `go test` with coverprofile
func GoTest(targets []string) {
	args := []string{"test"}
	args = append(args, targets...)
	args = append(args, "-coverprofile=cover.out")
	Run(goCommand, args...)
}

// GoGet for `go get`
func GoGet(packageName string) {
	Run(goCommand, "get", packageName)
}

// RunGoBin to run binary in gobin folder
func RunGoBin(name string, args ...string) {
	path := fmt.Sprintf("%s/%s/%s", build.Default.GOPATH, "bin", name)
	Run(path, args...)
}

// GoClean for `go clean`
func GoClean(cleanArgs ...string) {
	var args []string
	args = append(args, "clean")
	args = append(args, cleanArgs...)
	Run(goCommand, args...)
}
