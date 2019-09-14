package prebuilder

import (
	"io/ioutil"

	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

var (
	app        = typienv.App.MainPkg
	buildTool  = typienv.BuildTool.MainPkg
	dependency = "cmd/dependency"
)

// PreBuild process to build the typical project
func PreBuild(ctx *typictx.Context) (err error) {
	// path := ".typical"
	// if _, err := os.Stat(path); os.IsNotExist(err) {
	// 	os.Mkdir(path, os.ModePerm)
	// }
	root := "app"
	pkgs, filenames, _ := projectFiles(root)
	report, err := walker.Walk(filenames)
	if err != nil {
		return
	}
	// b, _ := json.MarshalIndent(report, "", "    ")
	// err = ioutil.WriteFile(path+"/walk_report.json", b, 0644)
	// if err != nil {
	// 	returnt
	// }
	pkg := "dependency"
	configuration := configuration(ctx)
	return runn.Execute(
		typienv.WriteEnvIfNotExist(ctx),
		genTestTargets(pkg, "test_targets.go", pkgs),
		genConstructors(pkg, "constructor.go", report),
		genConfiguration(pkg, "configuration.go", configuration),
		genSideEffects("main", "side_effects.go", ctx),
	)
}

func projectFiles(root string) (dirs []string, files []string, err error) {
	dirs = append(dirs, root)
	err = scanProjectFiles(root, &dirs, &files)
	return
}

func scanProjectFiles(root string, directories *[]string, files *[]string) (err error) {
	fileInfos, err := ioutil.ReadDir(root)
	if err != nil {
		return
	}
	for _, f := range fileInfos {
		if f.IsDir() {
			dirPath := root + "/" + f.Name()
			scanProjectFiles(dirPath, directories, files)
			*directories = append(*directories, dirPath)
		} else {
			*files = append(*files, root+"/"+f.Name())
		}
	}
	return
}

func genTestTargets(pkg, name string, testTargets []string) error {
	src := golang.NewSourceCode(pkg).AddTestTargets(testTargets...)
	target := dependency + "/" + name
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func genConstructors(pkg, name string, report *walker.Report) error {
	src := golang.NewSourceCode(pkg).
		AddConstructors(report.Autowires()...).
		AddMockTargets(report.Automocks()...)
	target := dependency + "/" + name
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func genConfiguration(pkg, name string, configuration ProjectConfiguration) error {
	src := golang.NewSourceCode(pkg).
		AddStruct(configuration.Struct).
		AddConstructorFunction(configuration.Constructors...)
	target := dependency + "/" + name
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func genSideEffects(pkg, name string, ctx *typictx.Context) error {
	appTarget := app + "/" + name
	devTarget := buildTool + "/" + name
	return runn.Execute(
		golang.NewSourceCode(pkg).AddImport(devToolSideEffects(ctx)...).Cook(devTarget),
		bash.GoImports(devTarget),
		golang.NewSourceCode(pkg).AddImport(appSideEffects(ctx)...).Cook(appTarget),
		bash.GoImports(appTarget),
	)
}
