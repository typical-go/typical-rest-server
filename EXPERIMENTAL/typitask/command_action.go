package typitask

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/util"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typigen"
	"gopkg.in/urfave/cli.v1"
)

func (t *TypicalTask) buildBinary(ctx *cli.Context) {
	typienv.GenerateAppEnvIfNotExist(t.Context)

	typigen.AppSideEffects(t.Context)

	binaryPath := typienv.BinaryPath(t.BinaryNameOrDefault())
	mainPackage := typienv.MainPackage(t.AppPkgOrDefault())

	log.Infof("Build the Binary for '%s' at '%s'", mainPackage, binaryPath)
	util.RunOrFatal(util.GoCommand(), "build", "-o", binaryPath, mainPackage)
}

func (t *TypicalTask) runBinary(ctx *cli.Context) {
	if !ctx.Bool("no-build") {
		t.buildBinary(ctx)
	}

	binaryPath := typienv.BinaryPath(t.BinaryNameOrDefault())
	log.Infof("Run the Binary '%s'", binaryPath)
	util.RunOrFatal(binaryPath, []string(ctx.Args())...)
}

func (t *TypicalTask) runTest(ctx *cli.Context) {
	log.Info("Run the Test")
	args := []string{"test"}
	args = append(args, t.AppModule.GetTestTargets()...)
	args = append(args, "-coverprofile=cover.out")
	util.RunOrFatal(util.GoCommand(), args...)
}

func (t *TypicalTask) releaseDistribution(ctx *cli.Context) {
	fmt.Println("Not implemented")
}

func (t *TypicalTask) generateMock(ctx *cli.Context) {
	util.RunOrFatal(util.GoCommand(), "get", "github.com/golang/mock/mockgen")
	mockPkg := t.MockPkgOrDefault()

	if ctx.Bool("new") {
		log.Infof("Clean mock package '%s'", mockPkg)
		os.RemoveAll(mockPkg)
	}

	for _, mockTarget := range t.AppModule.GetMockTargets() {
		dest := mockPkg + "/" + mockTarget[strings.LastIndex(mockTarget, "/")+1:]

		log.Infof("Generate mock for '%s' at '%s'", mockTarget, dest)
		util.RunOrFatal(util.GoBinary("mockgen"),
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

	log.Info("Trigger go clean")
	os.Setenv("GO111MODULE", "off") // NOTE:XXX: https://github.com/golang/go/issues/28680
	util.RunOrFatal(util.GoCommand(), "clean", "-x", "-testcache", "-modcachœœe")
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
