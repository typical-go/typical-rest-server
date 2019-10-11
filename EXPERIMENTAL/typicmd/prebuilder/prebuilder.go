package prebuilder

import (
	log "github.com/sirupsen/logrus"

	"time"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

// PreBuilder responsible to prebuild process
type PreBuilder struct {
	*typictx.Context
	*walker.ProjectFiles
	*walker.ContextFile
	Filenames []string
	Packages  []string
}

// TestTargets generate test target
func (p *PreBuilder) TestTargets() (err error) {
	defer elapsed("Generate TestTargets")()
	pkg := typienv.Dependency.Package
	src := golang.NewSourceCode(pkg)
	src.AddImport("", p.Root+"/typical")
	src.AddTestTargets(p.Packages...)
	target := dependency + "/test_targets.go"
	err = src.Cook(target)
	if err != nil {
		return
	}
	return bash.GoImports(target)

}

// Annotated to generate annotated
func (p *PreBuilder) Annotated() (err error) {
	defer elapsed("Generate Annotated")()
	pkg := typienv.Dependency.Package
	src := golang.NewSourceCode(pkg)
	for _, pkg := range p.Packages {
		src.AddImport("", p.Context.Root+"/"+pkg)
	}
	src.AddConstructors(p.ProjectFiles.Autowires()...)
	src.AddMockTargets(p.ProjectFiles.Automocks()...)
	target := dependency + "/annotateds.go"
	err = src.Cook(target)
	if err != nil {
		return
	}
	return bash.GoImports(target)
}

// Configuration to generate configuration
func (p *PreBuilder) Configuration() (err error) {
	defer elapsed("Generate Configuration")()
	conf := createConfiguration(p.Context)
	pkg := typienv.Dependency.Package
	src := golang.NewSourceCode(pkg).AddStruct(conf.Struct)
	src.AddImport("", "github.com/kelseyhightower/envconfig")
	for _, imp := range p.ContextFile.Imports {
		src.AddImport(imp.Name, imp.Path)
	}
	src.AddConstructors(conf.Constructors...)
	target := dependency + "/configurations.go"
	err = src.Cook(target)
	if err != nil {
		return
	}
	return bash.GoImports(target)
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		log.Debugf("%s took %v\n", what, time.Since(start))
	}
}
