package walker

import (
	"go/parser"
	"go/token"
	"strings"
)

// WalkContext to walk typical context
func WalkContext(ctxSrc string) (ctxFile *ContextFile, err error) {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, ctxSrc, nil, parser.ParseComments)
	if err != nil {
		return
	}
	ctxFile = new(ContextFile)
	for _, imp := range f.Imports {
		name := imp.Name.String()
		value := strings.Trim(imp.Path.Value, "\"")
		if name == "<nil>" {
			name = ""
		}
		ctxFile.AddImport(name, value)
	}
	return
}
