package typigen

import (
	"io/ioutil"

	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiast"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/golang"
)

var (
	app = typienv.App.MainPkg
	dev = typienv.BuildTool.MainPkg
)

// Generate all
func Generate(ctx *typictx.Context) (err error) {
	// path := ".typical"
	// if _, err := os.Stat(path); os.IsNotExist(err) {
	// 	os.Mkdir(path, os.ModePerm)
	// }
	root := "app"
	pkgs, filenames, _ := projectFiles(root)
	report, err := typiast.Walk(filenames)
	if err != nil {
		return
	}
	// b, _ := json.MarshalIndent(report, "", "    ")
	// err = ioutil.WriteFile(path+"/walk_report.json", b, 0644)
	// if err != nil {
	// 	returnt
	// }
	pkg := "main"
	configuration := configuration(ctx)
	return runn.Execute(
		typienv.WriteEnvIfNotExist(ctx),
		genTestTargets(pkg, "test_targets.go", pkgs),
		genConstructors(pkg, "constructor.go", report),
		genConfiguration(pkg, "configuration.go", configuration),
		genSideEffects(pkg, "side_effects.go", ctx),
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
	src := golang.NewSourceCode(pkg).
		AddTestTargets(testTargets...)
	appTarget := app + "/" + name
	devTarget := dev + "/" + name
	return runn.Execute(
		src.Cook(appTarget),
		bash.GoImports(appTarget),
		src.Cook(devTarget),
		bash.GoImports(devTarget),
	)
}

func genConstructors(pkg, name string, report *typiast.Report) error {
	src := golang.NewSourceCode(pkg).
		AddConstructors(report.Autowires()...).
		AddMockTargets(report.Automocks()...)
	appTarget := app + "/" + name
	devTarget := dev + "/" + name
	return runn.Execute(
		src.Cook(appTarget),
		bash.GoImports(appTarget),
		src.Cook(devTarget),
		bash.GoImports(devTarget),
	)
}

func genConfiguration(pkg, name string, configuration ProjectConfiguration) error {
	src := golang.NewSourceCode(pkg).
		AddStruct(configuration.Struct).
		AddConstructorFunction(configuration.Constructors...)
	appTarget := app + "/" + name
	devTarget := dev + "/" + name
	return runn.Execute(
		src.Cook(appTarget),
		bash.GoImports(appTarget),
		src.Cook(devTarget),
		bash.GoImports(devTarget),
	)
}

func genSideEffects(pkg, name string, ctx *typictx.Context) error {
	appTarget := app + "/" + name
	devTarget := dev + "/" + name
	return runn.Execute(
		golang.NewSourceCode(pkg).AddImport(devToolSideEffects(ctx)...).Cook(devTarget),
		bash.GoImports(devTarget),
		golang.NewSourceCode(pkg).AddImport(appSideEffects(ctx)...).Cook(appTarget),
		bash.GoImports(appTarget),
	)
}
