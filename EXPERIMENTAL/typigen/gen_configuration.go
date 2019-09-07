package typigen

import (
	"fmt"
	"reflect"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/golang"
)

// GenConfiguration generate configuration
func GenConfiguration(ctx *typictx.Context, srcCode *golang.SourceCode) (err error) {
	config, constructors := constructConfig(ctx)
	srcCode.AddStruct(config)
	srcCode.AddConstructorFunction(constructors...)
	return
}

func constructConfig(ctx *typictx.Context) (config golang.Struct, configConstructors []golang.Function) {
	structName := "Config"
	ptrStruct := "*" + structName
	config.Name = structName
	configConstructors = append(configConstructors, golang.Function{
		FuncParams:   map[string]string{},
		ReturnValues: []string{ptrStruct, "error"},
		FuncBody: fmt.Sprintf(`var cfg Config
err := envconfig.Process("", &cfg)
return &cfg, err`),
	})
	for _, acc := range ctx.ConfigAccessors() {
		key := acc.GetKey()
		typ := reflect.TypeOf(acc.GetConfigSpec())
		field := reflect.StructField{Name: key, Type: typ}
		function := golang.Function{
			FuncParams:   map[string]string{"cfg": ptrStruct},
			ReturnValues: []string{typ.String()},
			FuncBody:     fmt.Sprintf(`return cfg.%s`, key),
		}
		config.Fields = append(config.Fields, field)
		configConstructors = append(configConstructors, function)
	}
	return
}
