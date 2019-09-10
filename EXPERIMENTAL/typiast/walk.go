package typiast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
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
	paths, files, _ := projectFiles(appPath)
	report.Packages = paths
	fset := token.NewFileSet() // positions are relative to fset
	for _, filename := range files {
		if walkTarget(filename) {
			fmt.Println(filename)
			var f *ast.File
			f, err = parser.ParseFile(fset, filename, nil, parser.ParseComments)
			if err != nil {
				return
			}
			constructors, mock := parse(filename, f)
			report.Autowires = append(report.Autowires, constructors...)
			if mock {
				report.Automocks = append(report.Automocks, filename)
			}
		}
	}
	return
}

func walkTarget(filename string) bool {
	return strings.HasSuffix(filename, ".go") &&
		!strings.HasSuffix(filename, "_test.go")
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
