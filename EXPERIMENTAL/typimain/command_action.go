package typimain

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"gopkg.in/urfave/cli.v1"
)

func (t *TypicalTaskTool) buildBinary(ctx *cli.Context) {
	isGenerated, _ := generateNewEnviromentIfNotExist(t.Context)
	if isGenerated {
		log.Printf("Generate default enviroment at %s", envFile)
	}

	t.bundleAppSideEffects()

	binaryPath := typienv.BinaryPath(t.BinaryNameOrDefault())
	mainPackage := typienv.MainPackage(t.AppPkgOrDefault())

	log.Printf("Build the Binary for '%s' at '%s'", mainPackage, binaryPath)
	runOrFatal(goCommand(), "build", "-o", binaryPath, mainPackage)
}

func (t *TypicalTaskTool) runBinary(ctx *cli.Context) {
	if !ctx.Bool("no-build") {
		t.buildBinary(ctx)
	}

	binaryPath := typienv.BinaryPath(t.BinaryNameOrDefault())
	log.Printf("Run the Binary '%s'", binaryPath)
	runOrFatal(binaryPath, []string(ctx.Args())...)
}

func (t *TypicalTaskTool) runTest(ctx *cli.Context) {
	log.Println("Run the Test")
	args := []string{"test"}
	args = append(args, t.ArcheType.GetTestTargets()...)
	args = append(args, "-coverprofile=cover.out")
	runOrFatal(goCommand(), args...)
}

func (t *TypicalTaskTool) releaseDistribution(ctx *cli.Context) {
	fmt.Println("Not implemented")
}

func (t *TypicalTaskTool) generateMock(ctx *cli.Context) {
	runOrFatal(goCommand(), "get", "github.com/golang/mock/mockgen")
	mockPkg := t.MockPkgOrDefault()

	if ctx.Bool("new") {
		log.Printf("Clean mock package '%s'", mockPkg)
		os.RemoveAll(mockPkg)
	}

	for _, mockTarget := range t.ArcheType.GetMockTargets() {
		dest := mockPkg + "/" + mockTarget[strings.LastIndex(mockTarget, "/")+1:]

		log.Printf("Generate mock for '%s' at '%s'", mockTarget, dest)
		runOrFatal(goBinary("mockgen"),
			"-source", mockTarget,
			"-destination", dest,
			"-package", mockPkg)
	}
}

func (t *TypicalTaskTool) appPath(name string) string {
	return fmt.Sprintf("./%s/%s", t.AppPkgOrDefault(), name)
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
