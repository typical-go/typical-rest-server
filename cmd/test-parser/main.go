package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse src but stop after processing the imports.
	f, err := parser.ParseFile(fset, "typical/context.go", nil, 0)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Print the imports from the file's AST.
	// for _, s := range f.Imports {
	// 	fmt.Println(s.Path.Value)
	// }

	ctx, ok := f.Scope.Objects["Context"]
	if !ok {
		log.Fatal("Missing Context in typical/context.go")
		return
	}

	spec := ctx.Decl.(*ast.ValueSpec)
	composite := spec.Values[0].(*ast.CompositeLit)
	for _, expr := range composite.Elts {
		keyValue := expr.(*ast.KeyValueExpr)
		if fmt.Sprintf("%s", keyValue.Key) == "Modules" {
			composite2 := keyValue.Value.(*ast.CompositeLit)
			for _, expr2 := range composite2.Elts {
				switch expr2.(type) {
				case *ast.CompositeLit:
					module := parseToModule(expr2.(*ast.CompositeLit))
					fmt.Printf("%+v", module)
				case *ast.CallExpr:
				}

			}
		}
	}

	// fmt.Println(reflect.TypeOf(lit.Type.(*ast.SelectorExpr)))
}

func parseToModule(composite *ast.CompositeLit) (module typictx.Module) {
	for _, expr := range composite.Elts {
		keyValue := expr.(*ast.KeyValueExpr)

		fmt.Println(keyValue.Key)
	}
	return
}
