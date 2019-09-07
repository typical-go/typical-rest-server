package typigen

import (
	"fmt"
	"os"
	"reflect"

	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiparser"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/gosrc"

	log "github.com/sirupsen/logrus"
)

// TypicalGenerated to generate code in typical package
func TypicalGenerated(ctx *typictx.Context) (err error) {
	pkgName := "typical"
	filename := pkgName + "/generated.go"
	log.Infof("Typical Generated Code: %s", filename)
	mainConfig, configConstructors := constructConfig(ctx)
	projCtx, err := typiparser.Parse("app")
	if err != nil {
		return
	}
	recipe := gosrc.NewSourceCode(pkgName).
		AddStruct(mainConfig).
		AddConstructorFunction(configConstructors...).
		AddConstructors(projCtx.Autowires...).
		AddMockTargets(projCtx.Automocks...).
		AddTestTargets(projCtx.Packages...)

	if recipe.Blank() {
		os.Remove(filename)
		return
	}

	return runn.Execute(
		recipe.Cook(filename),
		bash.GoImports(pkgName),
	)
}

func constructConfig(ctx *typictx.Context) (mainConfig gosrc.Struct, configConstructors []gosrc.Function) {
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

	for _, acc := range ctx.ConfigAccessors() {
		field, function := configRecipe(acc, ptrMainConfigStruct)

		mainConfig.Fields = append(mainConfig.Fields, field)
		configConstructors = append(configConstructors, function)
	}

	return
}

func configRecipe(acc typictx.ConfigAccessor, ptrStruct string) (field reflect.StructField, function gosrc.Function) {
	key := acc.GetKey()
	typ := reflect.TypeOf(acc.GetConfigSpec())
	field = reflect.StructField{Name: key, Type: typ}
	function = gosrc.Function{
		FuncParams:   map[string]string{"cfg": ptrStruct},
		ReturnValues: []string{typ.String()},
		FuncBody:     fmt.Sprintf(`return cfg.%s`, key),
	}

	return
}
