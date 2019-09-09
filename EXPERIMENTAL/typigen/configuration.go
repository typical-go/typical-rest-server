package typigen

import (
	"fmt"
	"reflect"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/golang"
)

// ProjectConfiguration project configuration
type ProjectConfiguration struct {
	Struct       golang.Struct     `json:"struct"`
	Constructors []golang.Function `json:"constructors"`
}

func configuration(ctx *typictx.Context) (configuration ProjectConfiguration) {
	structName := "Config"
	ptrStruct := "*" + structName
	configuration.Struct.Name = structName
	configuration.Constructors = append(
		configuration.Constructors,
		golang.Function{
			FuncParams:   map[string]string{},
			ReturnValues: []string{ptrStruct, "error"},
			FuncBody: fmt.Sprintf(`var cfg Config
err := envconfig.Process("", &cfg)
return &cfg, err`),
		})
	for _, acc := range ctx.ConfigAccessors() {
		key := acc.GetKey()
		typ := reflect.TypeOf(acc.GetConfigSpec())
		configuration.Struct.Fields = append(
			configuration.Struct.Fields,
			reflect.StructField{Name: key, Type: typ},
		)
		configuration.Constructors = append(
			configuration.Constructors,
			golang.Function{
				FuncParams:   map[string]string{"cfg": ptrStruct},
				ReturnValues: []string{typ.String()},
				FuncBody:     fmt.Sprintf(`return cfg.%s`, key),
			},
		)
	}
	return
}
