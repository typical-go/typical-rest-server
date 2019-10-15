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

type prebuilder struct {
	*typictx.Context
	*walker.ProjectFiles
	*walker.ContextFile
	Filenames []string
	Packages  []string
}

func (p *prebuilder) Initiate(ctx *typictx.Context) (err error) {
	log.Debug("Scan project to get package and filenames")
	p.Context = ctx
	root := typienv.AppName
	p.Packages, p.Filenames, err = scanProject(root)
	if err != nil {
		return
	}
	log.Debug("Walk the project to get annotated or metadata")
	p.ProjectFiles, err = walker.WalkProject(p.Filenames)
	if err != nil {
		return
	}
	log.Debug("Walk the context file")
	p.ContextFile, err = walker.WalkContext(ctxPath)
	if err != nil {
		return
	}
	return
}

func (p *prebuilder) checkTestTargets() bool {
	return true
}

func (p *prebuilder) generateTestTargets() (err error) {
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

func (p *prebuilder) checkAnnotated() bool {
	return true
}

func (p *prebuilder) generateAnnotated() (err error) {
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

func (p *prebuilder) checkConfiguration() bool {
	return true
}

func (p *prebuilder) generateConfiguration() (err error) {
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
