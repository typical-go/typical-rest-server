package walker

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// WalkProject the source code to get autowire and automock
func WalkProject(filenames []string) (files *ProjectFiles, err error) {
	files = &ProjectFiles{}
	fset := token.NewFileSet() // positions are relative to fset
	for _, filename := range filenames {
		if isWalkTarget(filename) {
			var file ProjectFile
			file, err = parse(fset, filename)
			if err != nil {
				return
			}
			if !file.IsEmpty() {
				files.Add(file)
			}
		}
	}
	return
}

func isWalkTarget(filename string) bool {
	return strings.HasSuffix(filename, ".go") &&
		!strings.HasSuffix(filename, "_test.go")
}

func parse(fset *token.FileSet, filename string) (projFile ProjectFile, err error) {
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return
	}
	projFile.Name = filename
	for objName, obj := range f.Scope.Objects {
		switch obj.Decl.(type) {
		case *ast.FuncDecl:
			funcDecl := obj.Decl.(*ast.FuncDecl)
			var godoc string
			if funcDecl.Doc != nil {
				godoc = funcDecl.Doc.Text()
			}
			if isAutoWire(objName, godoc) {
				projFile.AddConstructor(fmt.Sprintf("%s.%s", f.Name, objName))
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
				projFile.Mock = isAutoMock(doc)
			}
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
