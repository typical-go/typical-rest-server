package typitask

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typigen"
	"gopkg.in/urfave/cli.v1"
)

func (t *TypicalTask) buildBinary(ctx *cli.Context) {
	typienv.GenerateAppEnvIfNotExist(t.Context)

	typigen.AppSideEffects(t.Context)

	binaryName := typienv.Binary(t.BinaryNameOrDefault())
	mainPackage := typienv.AppMainPackage()

	log.Infof("Build the Binary for '%s' at '%s'", mainPackage, binaryName)
	bash.GoBuild(binaryName, mainPackage)
}

func (t *TypicalTask) runBinary(ctx *cli.Context) {
	if !ctx.Bool("no-build") {
		t.buildBinary(ctx)
	}

	binaryPath := typienv.Binary(t.BinaryNameOrDefault())

	log.Infof("Run the Binary '%s'", binaryPath)
	bash.Run(binaryPath, []string(ctx.Args())...)
}

func (t *TypicalTask) runTest(ctx *cli.Context) {
	log.Info("Run the Test")
	bash.GoTest(t.AppModule.GetTestTargets())
}

func (t *TypicalTask) releaseDistribution(ctx *cli.Context) {

	t.runTest(ctx)
	t.generateReadme(ctx)

	goos := []string{"linux", "darwin"}
	goarch := []string{"amd64"}

	mainPackage := typienv.AppMainPackage()

	for _, os1 := range goos {
		os.Setenv("GOOS", os1)
		for _, arch := range goarch {
			// TODO: using ldflags
			binaryName := fmt.Sprintf("%s/%s_%s_%s", typienv.Release(), t.BinaryNameOrDefault(), os1, arch)
			os.Setenv("GOARCH", arch)

			log.Infof("Create release for %s/%s: %s", os1, arch, binaryName)
			bash.GoBuild(binaryName, mainPackage)
		}
	}

}

func (t *TypicalTask) generateMock(ctx *cli.Context) {
	bash.GoGet("github.com/golang/mock/mockgen")

	mockPkg := typienv.Mock()

	if ctx.Bool("new") {
		log.Infof("Clean mock package '%s'", mockPkg)
		os.RemoveAll(mockPkg)
	}

	for _, mockTarget := range t.AppModule.GetMockTargets() {
		dest := mockPkg + "/" + mockTarget[strings.LastIndex(mockTarget, "/")+1:]

		log.Infof("Generate mock for '%s' at '%s'", mockTarget, dest)
		bash.RunGoBin("mockgen",
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

	log.Infof("Generate ReadMe Document at '%s'", readmeFile)
	err = templ.Execute(file, Readme{
		Context: t.Context,
	})
	return nil
}

func (t *TypicalTask) cleanProject(ctx *cli.Context) {
	log.Info("Remove bin folder")
	os.RemoveAll(typienv.Bin())

	log.Info("Go clean")
	os.Setenv("GO111MODULE", "off") // NOTE:XXX: https://github.com/golang/go/issues/28680
	bash.GoClean("-x", "-testcache", "-modcachœœe")
}

func (t *TypicalTask) checkStatus(ctx *cli.Context) {
	statusReport := t.Context.CheckModuleStatus()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Module Name", "Status"})

	for moduleName, status := range statusReport {
		table.Append([]string{moduleName, status})
	}
	table.Render()

}
