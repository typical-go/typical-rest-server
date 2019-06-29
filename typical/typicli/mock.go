package typicli

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"

	"gopkg.in/urfave/cli.v1"
)

func mock(ctx *cli.Context) {
	runOrFatal(goCommand(), "get", "github.com/golang/mock/mockgen")

	// TODO: retrieve file based on pattern
	runOrFatal(goBinary("mockgen"), "-source=app/repository/book_repo.go", "-destination=mock/book_repo.go", "-package=mock")
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
