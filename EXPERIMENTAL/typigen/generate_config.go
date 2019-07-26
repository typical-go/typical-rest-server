package typigen

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"

	log "github.com/sirupsen/logrus"
)

// GenerateInitialization to generate initialization source code
func GenerateInitialization(ctx typictx.Context) (err error) {
	// TODO:
	log.Info("Generate typical initialization source code")

	generated := GeneratedModel{
		PackageName: "typical",
	}

	configStruct := StructPogo{Name: "Config"}

	generated.AddConstructors = append(generated.AddConstructors, FunctionPogo{
		FuncParams:   map[string]string{},
		ReturnValues: []string{"*Config", "error"},
		FuncBody: fmt.Sprintf(`var cfg Config
err := envconfig.Process("", &cfg)
return &cfg, err`),
	})

	for _, config := range ctx.Configs {
		typeConfig := reflect.TypeOf(config.Spec)
		configStruct.Fields = append(configStruct.Fields, reflect.StructField{
			Name: config.CamelPrefix(),
			Type: typeConfig,
		})
		generated.AddConstructors = append(generated.AddConstructors, FunctionPogo{
			FuncParams: map[string]string{
				"cfg": configStruct.Name,
			},
			ReturnValues: []string{typeConfig.String()},
			FuncBody:     fmt.Sprintf(`return cfg.%s`, config.CamelPrefix()),
		})

		for _, typ := range GetCompositionType(config.Spec) {
			generated.AddConstructors = append(generated.AddConstructors, FunctionPogo{
				FuncParams: map[string]string{
					"cfg": reflect.TypeOf(config.Spec).String(),
				},
				ReturnValues: []string{typ.String()},
				FuncBody:     `return cfg`,
			})
		}
	}

	generated.Structs = append(generated.Structs, configStruct)

	filename := "typical/generated.go"
	err = ioutil.WriteFile(filename, []byte(generated.String()), 0644)
	if err != nil {
		return
	}

	bash.GoImports(filename)
	return
}

func GetCompositionType(obj interface{}) (types []reflect.Type) {
	elm := reflect.TypeOf(obj).Elem()

	// loop through the struct's fields and set the map
	for i := 0; i < elm.NumField(); i++ {
		p := elm.Field(i)
		if p.Anonymous {
			types = append(types, p.Type)
		}
	}

	return
}
