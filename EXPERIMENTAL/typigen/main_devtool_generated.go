package typigen

import (
	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/golang"

	log "github.com/sirupsen/logrus"
)

// MainDevToolGenerated to generate code in typical package
func MainDevToolGenerated(ctx *typictx.Context) (err error) {
	filename := typienv.TypicalDevToolMainPackage() + "/generated.go"
	srcCode := golang.NewSourceCode("main")
	log.Infof("Generate recipe for Typical-Dev-Tool: %s", filename)
	return runn.Execute(
		GenDevToolSideEffects(ctx, srcCode),
		GenConfiguration(ctx, srcCode),
		GenDependecies(ctx, srcCode),
		srcCode.Cook(filename),
		bash.GoImports(filename),
	)
}
