package prebuilder

import (
	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

func generateApp(ctx *typictx.Context) error {
	pkg := "main"
	target := app + "/main.go"
	src := golang.NewSourceCode(pkg)
	src.Imports = appSideEffects(ctx)
	src.Put(mainApp())
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func mainApp() string {
	return `func main() {
	app := typicmd.NewTypicalApplication(typical.Context)
	err := app.Cli().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}`
}

func appSideEffects(ctx *typictx.Context) (imports golang.Imports) {
	for _, module := range ctx.Modules {
		for _, sideEffect := range module.SideEffects {
			if sideEffect.AppFlag {
				imports.Blank(sideEffect.Library)
			}
		}
	}
	imports.Blank(ctx.Root + "/" + dependency)
	imports.WithAlias("log", "github.com/sirupsen/logrus")
	imports.Add(ctx.Root + "/typical")
	return
}
