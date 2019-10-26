package prebuilder

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/pkg/utility/debugkit"
)

type testTarget struct {
	ContextImport string
	Packages      []string
}

func (g testTarget) generate(target string) (err error) {
	defer debugkit.ElapsedTime("Generate test target")()
	src := golang.NewSourceCode(typienv.Dependency.Package)
	src.AddImport("", g.ContextImport)
	src.AddTestTargets(g.Packages...)
	if err = src.Cook(target); err != nil {
		return
	}
	return bash.GoImports(target)
}
