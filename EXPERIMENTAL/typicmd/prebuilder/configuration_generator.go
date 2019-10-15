package prebuilder

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

// ConfigurationGenerator responsible to generate configuration
type ConfigurationGenerator struct {
	*typictx.Context
	*walker.ContextFile
}

// Generate the file
func (g *ConfigurationGenerator) Generate() (err error) {
	if g.check() {
		return g.generate()
	}
	return
}

func (g *ConfigurationGenerator) generate() (err error) {
	defer elapsed("Generate Configuration")()
	conf := createConfiguration(g.Context)
	pkg := typienv.Dependency.Package
	src := golang.NewSourceCode(pkg).AddStruct(conf.Struct)
	src.AddImport("", "github.com/kelseyhightower/envconfig")
	for _, imp := range g.ContextFile.Imports {
		src.AddImport(imp.Name, imp.Path)
	}
	src.AddConstructors(conf.Constructors...)
	target := dependency + "/configurations.go"
	err = src.Cook(target)
	if err != nil {
		return
	}
	return bash.GoImports(target)
	return
}

func (g *ConfigurationGenerator) check() bool {
	return true
}
