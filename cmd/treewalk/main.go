package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

func directoryFilter(f os.FileInfo) bool {
	return true
}

func main() {
	fset := token.NewFileSet() // positions are relative to fset
	paths := []string{"app"}

	for _, path := range paths {
		pkgs, err := parser.ParseDir(fset, path, directoryFilter, parser.ParseComments)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Print the imports from the file's AST.
		for pkgName, pkg := range pkgs {
			for fileName, file := range pkg.Files {
				for objName, obj := range file.Scope.Objects {
					switch obj.Decl.(type) {
					case *ast.FuncDecl:
						funcDecl := obj.Decl.(*ast.FuncDecl)
						var godoc string
						if funcDecl.Doc != nil {
							godoc = funcDecl.Doc.Text()
						}
						onFunctionFound(pkgName, objName, godoc)
					case *ast.TypeSpec:
						typeSpec := obj.Decl.(*ast.TypeSpec)
						switch typeSpec.Type.(type) {
						case *ast.StructType:
						case *ast.InterfaceType:
							var doc string
							if typeSpec.Doc != nil {
								doc = typeSpec.Doc.Text()
							}
							onInterfaceFound(fileName, objName, doc)
						}
					}
				}
			}
		}
	}
}

func onFunctionFound(pkgName, funcName, doc string) {
	fmt.Printf("Function: %s.%s: %s\n",
		pkgName, funcName, strings.TrimSpace(doc))
}

func onInterfaceFound(fileName, interfaceName, doc string) {
	fmt.Printf("Interface: %s in %s\n", interfaceName, fileName)
}
