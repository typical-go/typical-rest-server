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
	// 	return
	// }
	configuration := configuration(ctx)
	return runn.Execute(
		typienv.WriteEnvIfNotExist(ctx),
		appGenerated(ctx, configuration, pkgs, report),
		devToolGeneratead(ctx, configuration, pkgs, report),
	)
}

// func getCacheWalkReport() {
// }

func devToolGeneratead(ctx *typictx.Context, configuration ProjectConfiguration, testTargets []string, report *typiast.Report) error {
	pkg := "main"
	dir := typienv.TypicalDevToolMainPackage()
	return runn.Execute(
		genTestTargets(pkg, dir+"/test_targets.go", testTargets),
		genConstructors(pkg, dir+"/constructor.go", report),
		genConfiguration(pkg, dir+"/configuration.go", configuration),
		genSideEffects(pkg, dir+"/side_effects.go", devToolSideEffects(ctx)),
	)
}

func appGenerated(ctx *typictx.Context, configuration ProjectConfiguration, testTargets []string, report *typiast.Report) error {
	dir := typienv.AppMainPackage()
	pkg := "main"
	return runn.Execute(
		genTestTargets(pkg, dir+"/test_targets.go", testTargets),
		genConstructors(pkg, dir+"/constructor.go", report),
		genConfiguration(pkg, dir+"/configuration.go", configuration),
		genSideEffects(pkg, dir+"/side_effects.go", appSideEffects(ctx)),
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

func genTestTargets(pkg, target string, testTargets []string) error {
	src := golang.NewSourceCode(pkg).
		AddTestTargets(testTargets...)
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func genConstructors(pkg, target string, report *typiast.Report) error {
	src := golang.NewSourceCode(pkg).
		AddConstructors(report.Autowires()...).
		AddMockTargets(report.Automocks()...)
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func genConfiguration(pkg, target string, configuration ProjectConfiguration) error {
	src := golang.NewSourceCode(pkg).
		AddStruct(configuration.Struct).
		AddConstructorFunction(configuration.Constructors...)
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func genSideEffects(pkg, target string, sideEffects []golang.Import) error {
	src := golang.NewSourceCode(pkg).
		AddImport(sideEffects...)
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}
