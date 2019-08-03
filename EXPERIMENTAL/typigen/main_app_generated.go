package typigen

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/gosrc"
)

// MainAppGenerated to generate code in typical package
func MainAppGenerated(t typictx.Context) (err error) {
	filename := typienv.AppMainPackage() + "/generated.go"

	recipe := gosrc.Recipe{
		PackageName: "main",
	}

	for _, lib := range appSideEffects(t) {
		recipe.AddImport(gosrc.Import{Alias: "_", PackageName: lib})
	}

	if recipe.Blank() {
		os.Remove(filename)
		return
	}

	log.Infof("Generate recipe for App: %s", filename)
	return runn.Execute(
		recipe.Cook(filename),
		bash.GoFmt(filename),
	)
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
