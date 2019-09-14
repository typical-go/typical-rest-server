package prebuilder

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

func devToolSideEffects(ctx *typictx.Context) (imports []golang.Import) {
	imports = append(imports, golang.Import{
		Alias:       "_",
		PackageName: "github.com/typical-go/typical-rest-server/" + dependency,
	})
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
	imports = append(imports, golang.Import{
		Alias:       "_",
		PackageName: "github.com/typical-go/typical-rest-server/" + dependency,
	})
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
