package typimain

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
)

func (t *TypicalTaskTool) appPath(name string) string {
	return fmt.Sprintf("./%s/%s", t.AppPkgOrDefault(), name)
}

func goBinary(name string) string {
	return fmt.Sprintf("%s/%s/%s", build.Default.GOPATH, "bin", name)
}

func goCommand() string {
	return fmt.Sprintf("%s/bin/go", build.Default.GOROOT)
}

func runOrFatal(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func runOrFatalSilently(name string, args ...string) {
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
