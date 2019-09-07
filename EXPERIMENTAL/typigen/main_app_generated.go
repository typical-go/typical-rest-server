package typigen

import (
	log "github.com/sirupsen/logrus"

	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/golang"
)

// MainAppGenerated to generate code in typical package
func MainAppGenerated(ctx *typictx.Context) (err error) {
	filename := typienv.AppMainPackage() + "/generated.go"
	log.Infof("Generate recipe for App: %s", filename)
	srcCode := golang.NewSourceCode("main")
	return runn.Execute(
		GenAppSideEffects(ctx, srcCode),
		GenConfiguration(ctx, srcCode),
		GenDependecies(ctx, srcCode),
		srcCode.Cook(filename),
		bash.GoImports(filename),
	)
}
