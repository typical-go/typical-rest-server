package prebuilder

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/metadata"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/pkg/utility/debugkit"
)

// MockTargetGenerator responsible for generate mock_target
type MockTargetGenerator struct {
	Packages    []string
	Root        string
	MockTargets []string
}

// Generate the file
func (g *MockTargetGenerator) Generate() (updated bool, err error) {
	updated, err = metadata.Update("mock_target", g)
	if updated {
		err = g.generate()
	}
	return
}

func (g *MockTargetGenerator) generate() (err error) {
	defer debugkit.ElapsedTime("Generate Mock Target")()
	pkg := typienv.Dependency.Package
	src := golang.NewSourceCode(pkg)
	for _, pkg := range g.Packages {
		src.AddImport("", g.Root+"/"+pkg)
	}
	src.AddMockTargets(g.MockTargets...)
	target := dependency + "/mock_targets.go"
	err = src.Cook(target)
	if err != nil {
		return
	}
	return bash.GoImports(target)
}
