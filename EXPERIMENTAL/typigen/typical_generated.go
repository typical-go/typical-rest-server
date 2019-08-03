package typigen

import (
	"fmt"
	"os"
	"reflect"

	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typigen/generated"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiparser"

	log "github.com/sirupsen/logrus"
)

// TypicalGenerated to generate code in typical package
func TypicalGenerated(ctx typictx.Context) (err error) {
	// TODO: add typical folder in typienv
	packageName := "typical"
	filename := packageName + "/generated.go"
	log.Infof("Typical Generated Code: %s", filename)

	mainConfig, configConstructors := constructConfig(ctx)
	projCtx, err := typiparser.Parse("app")
	if err != nil {
		log.Fatal(err.Error())
	}

	recipe := generated.SourceRecipe{
		PackageName: packageName,
		Structs: []generated.StructPogo{
			mainConfig,
		},
	}
	recipe.AddConstructorPogos(configConstructors...)
	recipe.AddConstructors(projCtx.Autowires...)
	recipe.AddMockTargets(projCtx.Automocks...)
	recipe.AddTestTargets(projCtx.Packages...)

	if recipe.Blank() {
		os.Remove(filename)
		return
	}

	return runn.Execute(
		recipe.Cook(filename),
		bash.GoImports(packageName),
	)
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
