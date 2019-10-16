package prebuilder

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/metadata"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/pkg/utility/debugkit"
)

// TestTargetGenerator responsible to generate the test target
type TestTargetGenerator struct {
	Root     string
	Packages []string
}

// Generate the file
func (g *TestTargetGenerator) Generate() (updated bool, err error) {
	updated, err = metadata.Update("test_targets", g)
	if updated {
		err = g.generate()
	}
	return
}

func (g *TestTargetGenerator) generate() (err error) {
	defer debugkit.ElapsedTime("Generate TestTargets")()
	pkg := typienv.Dependency.Package
	src := golang.NewSourceCode(pkg)
	src.AddImport("", g.Root+"/typical")
	src.AddTestTargets(g.Packages...)
	target := dependency + "/test_targets.go"
	err = src.Cook(target)
	if err != nil {
		return
	}
	return bash.GoImports(target)
}
