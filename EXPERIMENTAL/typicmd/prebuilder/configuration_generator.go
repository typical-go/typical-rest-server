package prebuilder

import (
	"fmt"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/metadata"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/pkg/utility/debugkit"
)

// ConfigurationGenerator responsible to generate configuration
type ConfigurationGenerator struct {
	Configs []Config
	*walker.ContextFile
}

// Config model
type Config struct {
	Key string
	Typ string
}

// Generate the file
func (g *ConfigurationGenerator) Generate() (updated bool, err error) {
	updated, err = metadata.Update("configuration", g)
	if updated {
		err = g.generate()
	}
	return
}

func (g *ConfigurationGenerator) generate() (err error) {
	defer debugkit.ElapsedTime("Generate Configuration")()
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

func (g *ConfigurationGenerator) create() (model golang.Struct, constructors []string) {
	model.Name = "Config"
	model.Description = "for typical"
	constructors = append(constructors, g.configDef())
	for _, cfg := range g.Configs {
		model.AddField(cfg.Key, cfg.Typ)
		constructors = append(constructors, g.subConfigDef(cfg.Key, cfg.Typ))
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
