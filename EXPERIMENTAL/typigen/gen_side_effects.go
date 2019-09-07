package typigen

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/golang"
)

// GenAppSideEffects generate side effects for app
func GenAppSideEffects(ctx *typictx.Context, srcCode *golang.SourceCode) (err error) {
	for _, module := range ctx.Modules {
		for _, sideEffect := range module.SideEffects {
			if sideEffect.AppFlag {
				srcCode.AddImport(golang.Import{
					Alias:       "_",
					PackageName: sideEffect.Library,
				})
			}
		}
	}
	return
}

// GenDevToolSideEffects generate side effects for dev tool
func GenDevToolSideEffects(ctx *typictx.Context, srcCode *golang.SourceCode) (sideEffects []string) {
	for _, module := range ctx.Modules {
		for _, sideEffect := range module.SideEffects {
			if sideEffect.TypicalDevToolFlag {
				srcCode.AddImport(golang.Import{
					Alias:       "_",
					PackageName: sideEffect.Library,
				})
			}
		}
	}
	return
}
