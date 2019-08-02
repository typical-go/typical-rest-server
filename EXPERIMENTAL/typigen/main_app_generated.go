package typigen

import (
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/runn"
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
	}

	for _, lib := range appSideEffects(t) {
		recipe.AddImportPogo(generated.ImportPogo{Alias: "_", PackageName: lib})
	}

	return runn.Execute(
		recipe.Cook(filename), recipe.Cook(filename),
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
