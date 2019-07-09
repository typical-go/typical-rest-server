package typicli

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/typical-go/typical-rest-server/typical/typienv"
	"gopkg.in/urfave/cli.v1"
)

func (t *TypicalCli) updateTypical(ctx *cli.Context) {
	runOrFatal(goCommand(), "build", "-o", typienv.TypicalBinaryPath(), typienv.TypicalMainPackage())
}

func (t *TypicalCli) buildBinary(ctx *cli.Context) {
	GenerateNewEnviromentIfNotExist(t.Context)
	runOrFatal(goCommand(), "build", "-o",
		typienv.BinaryPath(t.TypiApp.BinaryName),
		typienv.MainPackage(t.TypiApp.ApplicationPkg))
}

func (t *TypicalCli) runBinary(ctx *cli.Context) {
	runOrFatal(typienv.BinaryPath(t.TypiApp.BinaryName), []string(ctx.Args())...)
}

func (t *TypicalCli) runTest(ctx *cli.Context) {
	args := []string{"test"}
	args = append(args, t.TypiApp.TestTargets...)
	args = append(args, "-coverprofile=cover.out")
	runOrFatal(goCommand(), args...)
}

func (t *TypicalCli) releaseDistribution(ctx *cli.Context) {
	fmt.Println("Not implemented")
}

func (t *TypicalCli) generateMock(ctx *cli.Context) {
	runOrFatal(goCommand(), "get", "github.com/golang/mock/mockgen")
	for _, mockTarget := range t.MockTargets {
		dest := t.MockPkg + "/" + mockTarget[strings.LastIndex(mockTarget, "/")+1:]
		runOrFatal(goBinary("mockgen"),
			"-source", mockTarget,
			"-destination", dest,
			"-package", t.MockPkg)
	}
}

func (t *TypicalCli) appPath(name string) string {
	return fmt.Sprintf("./%s/%s", t.ApplicationPkg, name)
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
