package prebuilder

import (
	"time"

	log "github.com/sirupsen/logrus"
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
	projFiles, err := walker.WalkProject(filenames)
	if err != nil {
		return
	}
	return runn.Execute(
		typienv.WriteEnvIfNotExist(ctx),
		prepareTestTargets(projPkgs),
		generateAnnotated(projFiles),
		generateConfiguration(configuration),
	)
}

func prepareTestTargets(projPkgs []string) (err error) {
	name := "test_targets"
	// err = writeCache(name, projPkgs)
	// if err != nil {
	// 	return
	// }
	return generateTestTargets(name, projPkgs)
}

func generateTestTargets(name string, testTargets []string) error {
	defer elapsed("Generate TestTargets")()
	pkg := typienv.Dependency.Package
	src := golang.NewSourceCode(pkg)
	src.AddTestTargets(testTargets...)
	target := dependency + "/" + name + ".go"
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func generateAnnotated(files *walker.ProjectFiles) error {
	defer elapsed("Generate Annotated")()
	pkg := typienv.Dependency.Package
	name := "annotateds.go"
	src := golang.NewSourceCode(pkg)
	src.AddConstructors(files.Autowires()...)
	src.AddMockTargets(files.Automocks()...)
	target := dependency + "/" + name
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func generateConfiguration(configuration Configuration) error {
	defer elapsed("Generate Configuration")()
	// TODO: try if manual import can improve goimport execution
	pkg := typienv.Dependency.Package
	name := "configurations.go"
	src := golang.NewSourceCode(pkg).
		AddStruct(configuration.Struct)
	src.AddImport("github.com/kelseyhightower/envconfig")
	src.AddConstructors(configuration.Constructors...)
	target := dependency + "/" + name
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		log.Infof("%s took %v\n", what, time.Since(start))
	}
}
