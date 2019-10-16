package prebuilder

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/metadata"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/pkg/utility/debugkit"
)

// AnnotatedGenerator responsible for generate annotated
type AnnotatedGenerator struct {
	*walker.ProjectFiles
	Packages []string
	Root     string
}

// Generate the file
func (g *AnnotatedGenerator) Generate() (err error) {
	updated, err := metadata.Update("annotated", g)
	if updated {
		return g.generate()
	}
	return
}

func (g *AnnotatedGenerator) generate() (err error) {
	defer debugkit.ElapsedTime("Generate Annotated")()
	pkg := typienv.Dependency.Package
	src := golang.NewSourceCode(pkg)
	for _, pkg := range g.Packages {
		src.AddImport("", g.Root+"/"+pkg)
	}
	src.AddConstructors(g.ProjectFiles.Autowires()...)
	src.AddMockTargets(g.ProjectFiles.Automocks()...)
	target := dependency + "/annotateds.go"
	err = src.Cook(target)
	if err != nil {
		return
	}
	return bash.GoImports(target)
}

func (g *AnnotatedGenerator) check() bool {
	return true
}
