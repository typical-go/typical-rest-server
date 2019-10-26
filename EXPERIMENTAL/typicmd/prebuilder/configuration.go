package prebuilder

import (
	"fmt"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/pkg/utility/debugkit"
)

type configuration struct {
	Configs       []config
	ConfigImports golang.Imports
}

type config struct {
	Key string
	Typ string
}

func (g configuration) generate(target string) (err error) {
	defer debugkit.ElapsedTime("Generate configuration")()
	model, contructors := g.create()
	pkg := typienv.Dependency.Package
	src := golang.NewSourceCode(pkg).AddStruct(model)
	src.Imports = g.ConfigImports
	src.AddConstructors(contructors...)
	if err = src.Cook(target); err != nil {
		return
	}
	return bash.GoImports(target)
}

func (g configuration) create() (model golang.Struct, constructors []string) {
	model.Name = "Config"
	model.Description = "for typical"
	constructors = append(constructors, g.configDef())
	for _, cfg := range g.Configs {
		model.AddField(cfg.Key, cfg.Typ)
		constructors = append(constructors, g.subConfigDef(cfg.Key, cfg.Typ))
	}
	return
}

func (g configuration) configDef() string {
	return `func() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return &cfg, err
}`
}

func (g *configuration) subConfigDef(name, typ string) string {
	return fmt.Sprintf(`func(cfg *Config) %s {
	return cfg.%s
}`, typ, name)
}
