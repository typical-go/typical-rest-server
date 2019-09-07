package typigen

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/golang"
)

func devToolSideEffects(ctx *typictx.Context) (imports []golang.Import) {
	for _, module := range ctx.Modules {
		for _, sideEffect := range module.SideEffects {
			if sideEffect.TypicalDevToolFlag {
				imports = append(imports, golang.Import{
					Alias:       "_",
					PackageName: sideEffect.Library,
				})
			}
		}
	}
	return
}

func appSideEffects(ctx *typictx.Context) (imports []golang.Import) {
	for _, module := range ctx.Modules {
		for _, sideEffect := range module.SideEffects {
			if sideEffect.AppFlag {
				imports = append(imports, golang.Import{
					Alias:       "_",
					PackageName: sideEffect.Library,
				})
			}
		}
	}
	return
}
