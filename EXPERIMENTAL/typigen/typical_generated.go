package typigen

import (
	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/golang"

	log "github.com/sirupsen/logrus"
)

// TypicalGenerated to generate code in typical package
func TypicalGenerated(ctx *typictx.Context) (err error) {
	pkgName := "typical"
	filename := pkgName + "/generated.go"
	log.Infof("Typical Generated Code: %s", filename)
	srcCode := golang.NewSourceCode(pkgName)
	return runn.Execute(
		GenConfiguration(ctx, srcCode),
		GenDependecies(ctx, srcCode),
		srcCode.Cook(filename),
		bash.GoImports(pkgName),
	)
}
