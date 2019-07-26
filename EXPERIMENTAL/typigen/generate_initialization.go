package typigen

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typigen/generated"

	log "github.com/sirupsen/logrus"
)

// GenerateInitialization to generate initialization source code
func GenerateInitialization(ctx typictx.Context) (err error) {
	// TODO:
	log.Info("Generate typical initialization source code")

	recipe := generated.SourceRecipe{
		PackageName: "typical",
	}

	configStructName := "Config"
	configStruct := generated.StructPogo{Name: configStructName}

	recipe.AddConstructors = append(recipe.AddConstructors, generated.FunctionPogo{
		FuncParams:   map[string]string{},
		ReturnValues: []string{"*" + configStructName, "error"},
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
		recipe.AddConstructors = append(recipe.AddConstructors, generated.FunctionPogo{
			FuncParams: map[string]string{
				"cfg": "*" + configStructName,
			},
			ReturnValues: []string{typeConfig.String()},
			FuncBody:     fmt.Sprintf(`return cfg.%s`, config.CamelPrefix()),
		})

		for _, typ := range GetCompositionType(config.Spec) {
			recipe.AddConstructors = append(recipe.AddConstructors, generated.FunctionPogo{
				FuncParams: map[string]string{
					"cfg": reflect.TypeOf(config.Spec).String(),
				},
				ReturnValues: []string{typ.String()},
				FuncBody:     `return cfg`,
			})
		}
	}

	recipe.Structs = append(recipe.Structs, configStruct)

	filename := "typical/generated.go"
	err = ioutil.WriteFile(filename, []byte(recipe.String()), 0644)
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
