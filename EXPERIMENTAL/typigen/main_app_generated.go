package typigen

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typigen/generated"
)

// MainAppGenerated to generate code in typical package
func MainAppGenerated(t typictx.Context) (err error) {
	filename := typienv.AppMainPackage() + "/generated.go"
	log.Infof("Main App Generated Code: %s", filename)

	recipe := generated.SourceRecipe{
		PackageName: "main",
		Imports:     make(map[string]string),
	}

	for _, lib := range appSideEffects(t) {
		recipe.Imports[lib] = "_"
	}

	err = ioutil.WriteFile(filename, []byte(recipe.String()), 0644)
	if err != nil {
		return
	}

	bash.GoFmt(filename)
	return
}

func appSideEffects(t typictx.Context) (sideEffects []string) {
	for _, module := range t.Modules {
		for _, sideEffect := range module.SideEffects {
			if sideEffect.AppFlag {
				sideEffects = append(sideEffects, sideEffect.Library)
			}
		}
	}
	return
}
