package prebuilder

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/pkg/utility/debugkit"
)

type constructor struct {
	ApplicationImports golang.Imports
	Constructors       []string
}

func (g constructor) generate(target string) (err error) {
	defer debugkit.ElapsedTime("Generate constructor")()
	src := golang.NewSourceCode(typienv.Dependency.Package)
	src.Imports = g.ApplicationImports
	src.AddConstructors(g.Constructors...)
	if err = src.Cook(target); err != nil {
		return
	}
	return bash.GoImports(target)
}
