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

// TypicalGenerated to generate code in typical package
func TypicalGenerated(ctx typictx.Context) (err error) {
	log.Info("Generate typical initialization source code")

	mainConfig, configConstructors := constructConfig(ctx)

	recipe := generated.SourceRecipe{
		PackageName: "typical",
		Structs: []generated.StructPogo{
			mainConfig,
		},
	}
	recipe.AddConstructors = append(recipe.AddConstructors, configConstructors...)

	// TODO: add typical folder in typienv
	filename := "typical/generated.go"
	err = ioutil.WriteFile(filename, []byte(recipe.String()), 0644)
	if err != nil {
		return
	}

	bash.GoImports(filename)
	return
}

func constructConfig(ctx typictx.Context) (mainConfig generated.StructPogo, configConstructors []generated.FunctionPogo) {
	mainConfigStruct := "Config"
	ptrMainConfigStruct := "*" + mainConfigStruct

	mainConfig.Name = mainConfigStruct

	configConstructors = append(configConstructors, generated.FunctionPogo{
		FuncParams:   map[string]string{},
		ReturnValues: []string{ptrMainConfigStruct, "error"},
		FuncBody: fmt.Sprintf(`var cfg Config
err := envconfig.Process("", &cfg)
return &cfg, err`),
	})

	for _, config := range ctx.Configs {
		typeConfig := reflect.TypeOf(config.Spec)
		mainConfig.Fields = append(mainConfig.Fields, reflect.StructField{
			Name: config.CamelPrefix(),
			Type: typeConfig,
		})
		configConstructors = append(configConstructors, generated.FunctionPogo{
			FuncParams:   map[string]string{"cfg": ptrMainConfigStruct},
			ReturnValues: []string{typeConfig.String()},
			FuncBody:     fmt.Sprintf(`return cfg.%s`, config.CamelPrefix()),
		})
	}

	return

}
