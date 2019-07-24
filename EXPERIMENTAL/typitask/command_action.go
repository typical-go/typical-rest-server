package typitask

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typigen"
)

// BuildBinary for typical application
func BuildBinary(ctx typictx.ActionContext) error {
	typienv.GenerateAppEnvIfNotExist(ctx.Typical)
	typigen.AppSideEffects(ctx.Typical)

	binaryName := typienv.Binary(ctx.Typical.BinaryNameOrDefault())
	mainPackage := typienv.AppMainPackage()

	log.Infof("Build the Binary for '%s' at '%s'", mainPackage, binaryName)
	bash.GoBuild(binaryName, mainPackage)

	return nil
}

// RunBinary for run typical binary
func RunBinary(ctx typictx.ActionContext) error {
	if !ctx.Cli.Bool("no-build") {
		BuildBinary(ctx)
	}

	binaryPath := typienv.Binary(ctx.Typical.BinaryNameOrDefault())

	log.Infof("Run the Binary '%s'", binaryPath)
	bash.Run(binaryPath, []string(ctx.Cli.Args())...)
	return nil
}

func RunTest(ctx typictx.ActionContext) error {
	log.Info("Run the Test")
	bash.GoTest(ctx.Typical.AppModule.GetTestTargets())
	return nil
}

func ReleaseDistribution(ctx typictx.ActionContext) error {
	RunTest(ctx)
	GenerateReadme(ctx)

	goos := []string{"linux", "darwin"}
	goarch := []string{"amd64"}
	mainPackage := typienv.AppMainPackage()

	for _, os1 := range goos {
		os.Setenv("GOOS", os1)
		for _, arch := range goarch {
			// TODO: using ldflags
			binaryName := fmt.Sprintf("%s/%s_%s_%s",
				typienv.Release(), ctx.Typical.BinaryNameOrDefault(), os1, arch)
			os.Setenv("GOARCH", arch)

			log.Infof("Create release for %s/%s: %s", os1, arch, binaryName)
			bash.GoBuild(binaryName, mainPackage)
		}
	}
	return nil
}

func GenerateMock(ctx typictx.ActionContext) error {
	bash.GoGet("github.com/golang/mock/mockgen")

	mockPkg := typienv.Mock()

	if ctx.Cli.Bool("new") {
		log.Infof("Clean mock package '%s'", mockPkg)
		os.RemoveAll(mockPkg)
	}

	for _, mockTarget := range ctx.Typical.AppModule.GetMockTargets() {
		dest := mockPkg + "/" + mockTarget[strings.LastIndex(mockTarget, "/")+1:]

		log.Infof("Generate mock for '%s' at '%s'", mockTarget, dest)
		bash.RunGoBin("mockgen",
			"-source", mockTarget,
			"-destination", dest,
			"-package", mockPkg)
	}
	return nil
}

// GenerateReadme for generate typical applical readme
func GenerateReadme(ctx typictx.ActionContext) (err error) {
	readmeFile := ctx.Typical.ReadmeFileOrDefault()
	readmeTemplate := ctx.Typical.ReadmeTemplateOrDefault()

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
		Context: ctx.Typical,
	})
	return nil
}

func CleanProject(ctx typictx.ActionContext) error {
	log.Info("Remove bin folder")
	os.RemoveAll(typienv.Bin())

	log.Info("Go clean")
	os.Setenv("GO111MODULE", "off") // NOTE:XXX: https://github.com/golang/go/issues/28680
	bash.GoClean("-x", "-testcache", "-modcachœœe")
	return nil
}

func CheckStatus(ctx typictx.ActionContext) error {
	statusReport := ctx.Typical.CheckModuleStatus()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Module Name", "Status"})

	for moduleName, status := range statusReport {
		table.Append([]string{moduleName, status})
	}
	table.Render()
	return nil
}
