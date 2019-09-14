package prebuilder

import (
	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

func generateBuildTool(ctx *typictx.Context) error {
	pkg := "main"
	target := buildTool + "/main.go"
	src := golang.NewSourceCode(pkg)
	src.Imports = buildToolImports(ctx)
	src.Put(mainBuildTool())
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func mainBuildTool() string {
	return `func main() {
	buildTool := typicmd.NewTypicalBuildTool(typical.Context)
	err := buildTool.Cli().Run(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}`
}

func buildToolImports(ctx *typictx.Context) (imports golang.Imports) {
	for _, module := range ctx.Modules {
		for _, sideEffect := range module.SideEffects {
			if sideEffect.TypicalDevToolFlag {
				imports.Blank(sideEffect.Library)
			}
		}
	}
	imports.Blank(ctx.Root + "/" + dependency)
	imports.WithAlias("log", "github.com/sirupsen/logrus")
	imports.Add(ctx.Root + "/typical")
	return
}
