package prebuilder

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/pkg/utility/debugkit"
)

type mockTarget struct {
	ApplicationImports golang.Imports
	MockTargets        []string
}

func (g mockTarget) generate(target string) (err error) {
	defer debugkit.ElapsedTime("Generate mock target")()
	src := golang.NewSourceCode(typienv.Dependency.Package)
	src.Imports = g.ApplicationImports
	src.AddMockTargets(g.MockTargets...)
	if err = src.Cook(target); err != nil {
		return
	}
	return bash.GoImports(target)
}
