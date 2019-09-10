package typiast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
)

// Report of goAst when walk the project
type Report struct {
	Packages  []string `json:"packages"`
	Autowires []string `json:"autowires"`
	Automocks []string `json:"automocks"`
}

// Walk the source code to get autowire and automock
func Walk(appPath string) (report *Report, err error) {
	report = &Report{}
	paths := []string{appPath}
	testTargets := make(map[string]struct{})
	allDir(appPath, &paths)
	fmt.Println(paths)
	fset := token.NewFileSet() // positions are relative to fset
	for _, path := range paths {
		var pkgs map[string]*ast.Package
		pkgs, err = parser.ParseDir(fset, path, directoryFilter, parser.ParseComments)
		if err != nil {
			return
		}
		// Print the imports from the file's AST.
		for _, pkg := range pkgs {
			testTargets[path] = struct{}{}
			for fileName, file := range pkg.Files {
				constructors, mock := parse(fileName, file)
				report.Autowires = append(report.Autowires, constructors...)
				if mock {
					report.Automocks = append(report.Automocks, fileName)
				}

			}
		}
	}
	for key := range testTargets {
		report.Packages = append(report.Packages, key)
	}
	return
}

func parse(filename string, file *ast.File) (constructors []string, mock bool) {
	for objName, obj := range file.Scope.Objects {
		switch obj.Decl.(type) {
		case *ast.FuncDecl:
			funcDecl := obj.Decl.(*ast.FuncDecl)
			var godoc string
			if funcDecl.Doc != nil {
				godoc = funcDecl.Doc.Text()
			}
			if isAutoWire(objName, godoc) {
				constructors = append(constructors, fmt.Sprintf("%s.%s", file.Name, objName))
			}
		case *ast.TypeSpec:
			typeSpec := obj.Decl.(*ast.TypeSpec)
			switch typeSpec.Type.(type) {
			case *ast.StructType:
			case *ast.InterfaceType:
				var doc string
				if typeSpec.Doc != nil {
					doc = typeSpec.Doc.Text()
				}
				mock = isAutoMock(doc)
			}
		}
	}
	return
}

func allDir(path string, directories *[]string) (err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, f := range files {
		if f.IsDir() {
			dirPath := path + "/" + f.Name()
			allDir(dirPath, directories)
			*directories = append(*directories, dirPath)
		}
	}
	return
}

func directoryFilter(f os.FileInfo) bool {
	return true
}

func isAutoWire(funcName, doc string) bool {
	tags := ParseDocTag(doc)
	if strings.HasPrefix(funcName, "New") {
		return !tags.Contain("nowire")
	}
	return tags.Contain("autowire")
}

func isAutoMock(doc string) bool {
	tags := ParseDocTag(doc)
	return !tags.Contain("nomock")
}
