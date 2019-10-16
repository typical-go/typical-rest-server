package prebuilder

import (
	"fmt"
	"reflect"

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
	model, contructors := g.create()
	pkg := typienv.Dependency.Package
	src := golang.NewSourceCode(pkg).AddStruct(model)
	src.AddImport("", "github.com/kelseyhightower/envconfig")
	for _, imp := range g.ContextFile.Imports {
		src.AddImport(imp.Name, imp.Path)
	}
	src.AddConstructors(contructors...)
	target := dependency + "/configurations.go"
	err = src.Cook(target)
	if err != nil {
		return
	}
	return bash.GoImports(target)
}

func (g *ConfigurationGenerator) check() bool {
	return true
}

func (g *ConfigurationGenerator) create() (model golang.Struct, constructors []string) {
	structName := "Config"
	model.Name = structName
	constructors = append(constructors, g.configDef())
	for _, acc := range g.ConfigAccessors() {
		key := acc.GetKey()
		typ := reflect.TypeOf(acc.GetConfigSpec()).String()
		model.AddField(key, typ)
		constructors = append(constructors, g.subConfigDef(key, typ))

	}
	return
}

func (g *ConfigurationGenerator) configDef() string {
	return `func() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return &cfg, err
}`
}

func (g *ConfigurationGenerator) subConfigDef(name, typ string) string {
	return fmt.Sprintf(`func(cfg *Config) %s {
	return cfg.%s
}`, typ, name)
}
