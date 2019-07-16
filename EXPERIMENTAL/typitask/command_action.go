package typitask

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"gopkg.in/urfave/cli.v1"
)

func (t *TypicalTask) buildBinary(ctx *cli.Context) {
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

func (t *TypicalTask) runBinary(ctx *cli.Context) {
	if !ctx.Bool("no-build") {
		t.buildBinary(ctx)
	}

	binaryPath := typienv.BinaryPath(t.BinaryNameOrDefault())
	log.Printf("Run the Binary '%s'", binaryPath)
	runOrFatal(binaryPath, []string(ctx.Args())...)
}

func (t *TypicalTask) runTest(ctx *cli.Context) {
	log.Println("Run the Test")
	args := []string{"test"}
	args = append(args, t.ArcheType.GetTestTargets()...)
	args = append(args, "-coverprofile=cover.out")
	runOrFatal(goCommand(), args...)
}

func (t *TypicalTask) releaseDistribution(ctx *cli.Context) {
	fmt.Println("Not implemented")
}

func (t *TypicalTask) generateMock(ctx *cli.Context) {
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

func (t *TypicalTask) generateReadme(ctx *cli.Context) (err error) {
	readmeFile := t.ReadmeFileOrDefault()
	readmeTemplate := t.ReadmeTemplateOrDefault()

	templ, err := template.New("readme").Parse(readmeTemplate)

	if err != nil {
		return
	}

	file, err := os.Create(readmeFile)
	if err != nil {
		return
	}

	log.Printf("Generate ReadMe Document at '%s'", readmeFile)
	err = templ.Execute(file, Readme{
		Context: t.Context,
	})
	return nil
}

func (t *TypicalTask) cleanProject(ctx *cli.Context) {
	log.Println("Remove bin folder")
	os.RemoveAll(typienv.Bin())

	log.Println("Trigger go clean")
	os.Setenv("GO111MODULE", "off") // NOTE:XXX: https://github.com/golang/go/issues/28680
	runOrFatal(goCommand(), "clean", "-x", "-testcache", "-modcache")
}
