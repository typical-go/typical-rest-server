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
	log.Println("Update the typical")
	runOrFatal(goCommand(), "build", "-o", typienv.TypicalBinaryPath(), typienv.TypicalMainPackage())
}

func (t *TypicalCli) buildBinary(ctx *cli.Context) {
	isGenerated, _ := generateNewEnviromentIfNotExist(t.Context)
	if isGenerated {
		log.Printf("Generate default enviroment at %s", envFile)
	}

	binaryPath := typienv.BinaryPath(t.TypiApp.BinaryNameOrDefault())
	mainPackage := typienv.MainPackage(t.TypiApp.ApplicationPkgOrDefault())

	log.Printf("Build the Binary for '%s' at '%s'", mainPackage, binaryPath)
	runOrFatal(goCommand(), "build", "-o", binaryPath, mainPackage)
}

func (t *TypicalCli) runBinary(ctx *cli.Context) {
	if !ctx.Bool("no-build") {
		t.buildBinary(ctx)
	}

	binaryPath := typienv.BinaryPath(t.TypiApp.BinaryNameOrDefault())
	log.Printf("Run the Binary '%s'", binaryPath)
	runOrFatal(binaryPath, []string(ctx.Args())...)
}

func (t *TypicalCli) runTest(ctx *cli.Context) {
	log.Println("Run the Test")
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

	mockPkg := t.MockPkgOrDefault()

	log.Printf("Clean mock package '%s'", mockPkg)
	os.RemoveAll(mockPkg)
	for _, mockTarget := range t.MockTargets {
		dest := mockPkg + "/" + mockTarget[strings.LastIndex(mockTarget, "/")+1:]

		log.Printf("Generate mock for '%s' at '%s'", mockTarget, dest)
		runOrFatal(goBinary("mockgen"),
			"-source", mockTarget,
			"-destination", dest,
			"-package", mockPkg)
	}
}

func (t *TypicalCli) appPath(name string) string {
	return fmt.Sprintf("./%s/%s", t.ApplicationPkgOrDefault(), name)
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
