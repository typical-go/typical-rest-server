package prebuilder

import (
	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

func generateDependency(ctx *typictx.Context) error {
	root := typienv.AppName
	pkg := typienv.Dependency.Package
	pkgs, filenames, _ := projectFiles(root)
	configuration := createConfiguration(ctx)
	report, err := walker.Walk(filenames)
	if err != nil {
		return err
	}
	return runn.Execute(
		genTestTargets(pkg, "test_targets.go", pkgs),
		genConstructors(pkg, "constructors.go", report),
		genConfiguration(pkg, "configurations.go", configuration),
	)
}

func genTestTargets(pkg, name string, testTargets []string) error {
	src := golang.NewSourceCode(pkg)
	src.AddTestTargets(testTargets...)
	target := dependency + "/" + name
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func genConstructors(pkg, name string, report *walker.Report) error {
	src := golang.NewSourceCode(pkg)
	src.AddConstructors(report.Autowires()...)
	src.AddMockTargets(report.Automocks()...)
	target := dependency + "/" + name
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func genConfiguration(pkg, name string, configuration Configuration) error {
	src := golang.NewSourceCode(pkg).
		AddStruct(configuration.Struct)
	src.AddConstructors(configuration.Constructors...)
	target := dependency + "/" + name
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}
