package prebuilder

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/metadata"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/pkg/utility/debugkit"
)

// ConstructorGenerator responsible for generate annotated
type ConstructorGenerator struct {
	Packages     []string
	Root         string
	Constructors []string
}

// Generate the file
func (g *ConstructorGenerator) Generate() (updated bool, err error) {
	updated, err = metadata.Update("constructor", g)
	if updated {
		err = g.generate()
	}
	return
}

func (g *ConstructorGenerator) generate() (err error) {
	defer debugkit.ElapsedTime("Generate Constructors")()
	pkg := typienv.Dependency.Package
	src := golang.NewSourceCode(pkg)
	for _, pkg := range g.Packages {
		src.AddImport("", g.Root+"/"+pkg)
	}
	src.AddConstructors(g.Constructors...)
	target := dependency + "/constructors.go"
	err = src.Cook(target)
	if err != nil {
		return
	}
	return bash.GoImports(target)
}
