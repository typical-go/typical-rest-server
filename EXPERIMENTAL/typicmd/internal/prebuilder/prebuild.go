package prebuilder

import (
	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

var (
	app        = typienv.App.SrcPath
	buildTool  = typienv.BuildTool.SrcPath
	dependency = typienv.Dependency.SrcPath
)

// PreBuild process to build the typical project
func PreBuild(ctx *typictx.Context) (err error) {
	root := typienv.AppName
	projPkgs, filenames, _ := projectFiles(root)
	configuration := createConfiguration(ctx)
	report, err := walker.Walk(filenames)
	if err != nil {
		return err
	}
	return runn.Execute(
		typienv.WriteEnvIfNotExist(ctx),
		generateTestTargets(projPkgs),
		generateAnnotated(report),
		generateConfiguration(configuration),
	)
}

func generateTestTargets(testTargets []string) error {
	pkg := typienv.Dependency.Package
	name := "test_targets.go"
	src := golang.NewSourceCode(pkg)
	src.AddTestTargets(testTargets...)
	target := dependency + "/" + name
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func generateAnnotated(report *walker.Report) error {
	pkg := typienv.Dependency.Package
	name := "annotateds.go"
	src := golang.NewSourceCode(pkg)
	src.AddConstructors(report.Autowires()...)
	src.AddMockTargets(report.Automocks()...)
	target := dependency + "/" + name
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func generateConfiguration(configuration Configuration) error {
	pkg := typienv.Dependency.Package
	name := "configurations.go"
	src := golang.NewSourceCode(pkg).
		AddStruct(configuration.Struct)
	src.AddConstructors(configuration.Constructors...)
	target := dependency + "/" + name
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}
