package walker

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// Walk the source code to get autowire and automock
func Walk(filenames []string) (f *Files, err error) {
	f = &Files{}
	fset := token.NewFileSet() // positions are relative to fset
	for _, filename := range filenames {
		if walkTarget(filename) {
			var file File
			file, err = parse(fset, filename)
			if err != nil {
				return
			}
			if !file.IsEmpty() {
				f.Add(file)
			}
		}
	}
	return
}

func walkTarget(filename string) bool {
	return strings.HasSuffix(filename, ".go") &&
		!strings.HasSuffix(filename, "_test.go")
}

func parse(fset *token.FileSet, filename string) (file File, err error) {
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return
	}
	file.Name = filename
	for objName, obj := range f.Scope.Objects {
		switch obj.Decl.(type) {
		case *ast.FuncDecl:
			funcDecl := obj.Decl.(*ast.FuncDecl)
			var godoc string
			if funcDecl.Doc != nil {
				godoc = funcDecl.Doc.Text()
			}
			if isAutoWire(objName, godoc) {
				file.AddConstructor(fmt.Sprintf("%s.%s", f.Name, objName))
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
				file.Mock = isAutoMock(doc)
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
