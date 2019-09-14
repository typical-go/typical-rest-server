package prebuilder

import (
	"fmt"
	"reflect"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

// Configuration project configuration
type Configuration struct {
	Struct       golang.Struct `json:"struct"`
	Constructors []string      `json:"constructors"`
}

// AddConstructor to add constructor to project configuration
func (c *Configuration) AddConstructor(constructor string) {
	c.Constructors = append(c.Constructors, constructor)
}

func createConfiguration(ctx *typictx.Context) (cfg Configuration) {
	structName := "Config"
	cfg.Struct.Name = structName
	cfg.AddConstructor(configDefinition())
	for _, acc := range ctx.ConfigAccessors() {
		key := acc.GetKey()
		typ := reflect.TypeOf(acc.GetConfigSpec())
		cfg.Struct.AddField(reflect.StructField{Name: key, Type: typ})
		cfg.AddConstructor(subConfigDefinition(key, typ.String()))
	}
	return
}

func configDefinition() string {
	return `func() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return &cfg, err
}`
}

func subConfigDefinition(name, typ string) string {
	return fmt.Sprintf(`func(cfg *Config) %s {
	return cfg.%s
}`, typ, name)
}
