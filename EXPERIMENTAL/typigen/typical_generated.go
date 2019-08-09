package typigen

import (
	"fmt"
	"os"
	"reflect"

	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiparser"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/gosrc"

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

	recipe := gosrc.Recipe{
		PackageName: packageName,
		Structs: []gosrc.Struct{
			mainConfig,
		},
	}
	recipe.AddConstructorFunction(configConstructors...)
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

func constructConfig(ctx typictx.Context) (mainConfig gosrc.Struct, configConstructors []gosrc.Function) {
	mainConfigStruct := "Config"
	ptrMainConfigStruct := "*" + mainConfigStruct

	mainConfig.Name = mainConfigStruct

	configConstructors = append(configConstructors, gosrc.Function{
		FuncParams:   map[string]string{},
		ReturnValues: []string{ptrMainConfigStruct, "error"},
		FuncBody: fmt.Sprintf(`var cfg Config
err := envconfig.Process("", &cfg)
return &cfg, err`),
	})

	for _, config := range ctx.Configurations {
		typeConfig := reflect.TypeOf(config.Spec)
		mainConfig.Fields = append(mainConfig.Fields, reflect.StructField{
			Name: config.CamelPrefix(),
			Type: typeConfig,
		})
		configConstructors = append(configConstructors, gosrc.Function{
			FuncParams:   map[string]string{"cfg": ptrMainConfigStruct},
			ReturnValues: []string{typeConfig.String()},
			FuncBody:     fmt.Sprintf(`return cfg.%s`, config.CamelPrefix()),
		})
	}

	return
}
