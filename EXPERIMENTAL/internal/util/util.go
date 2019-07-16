package util // TODO: rename with more descriptive name


import (
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
)

// GoBinary return absolute path of binary in GOBIN
// TODO: create run wrapper
func GoBinary(name string) string {
	return fmt.Sprintf("%s/%s/%s", build.Default.GOPATH, "bin", name)
}

// GoCommand return absoulte path of go command
// TODO: create run wrapper
func GoCommand() string {
	return fmt.Sprintf("%s/bin/go", build.Default.GOROOT)
}

// RunOrFatal to execute bash or exit the application
func RunOrFatal(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// RunOrFatalSilently same with RunOrFatal without print any output
func RunOrFatalSilently(name string, args ...string) {
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
