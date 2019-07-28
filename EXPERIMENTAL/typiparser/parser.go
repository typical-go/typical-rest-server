package typiparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// Parse the source code to get autowire and automock
func Parse(appPath string) (projCtx ProjectContext, err error) {
	var paths []string
	testTargets := make(map[string]struct{})

	AllDirectories(appPath, &paths)
	fset := token.NewFileSet() // positions are relative to fset

	for _, path := range paths {
		var pkgs map[string]*ast.Package
		pkgs, err = parser.ParseDir(fset, path, directoryFilter, parser.ParseComments)
		if err != nil {
			return
		}

		// Print the imports from the file's AST.
		for pkgName, pkg := range pkgs {
			testTargets[path] = struct{}{}
			for fileName, file := range pkg.Files {
				for objName, obj := range file.Scope.Objects {
					switch obj.Decl.(type) {
					case *ast.FuncDecl:
						funcDecl := obj.Decl.(*ast.FuncDecl)
						var godoc string
						if funcDecl.Doc != nil {
							godoc = funcDecl.Doc.Text()
						}
						if isAutoWire(objName, godoc) {
							projCtx.Autowires = append(projCtx.Autowires, fmt.Sprintf("%s.%s", pkgName, objName))
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
							if isAutoMock(doc) {
								projCtx.Automocks = append(projCtx.Automocks, fileName)
							}
						}
					}
				}
			}
		}
	}

	for key := range testTargets {
		projCtx.Packages = append(projCtx.Packages, key)
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
