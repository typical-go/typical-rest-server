package typigen

import (
	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiast"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/golang"
)

// Generate all
func Generate(ctx *typictx.Context) (err error) {
	out, err := typiast.Walk("app")
	if err != nil {
		return
	}
	configuration := configuration(ctx)
	appOut := typienv.AppMainPackage() + "/generated.go"
	devOut := typienv.TypicalDevToolMainPackage() + "/generated.go"
	return runn.Execute(
		typienv.WriteEnvIfNotExist(ctx),
		appSource(ctx, configuration, out).Cook(appOut),
		devToolSource(ctx, configuration, out).Cook(devOut),
		bash.GoImports(appOut),
		bash.GoImports(devOut),
	)
}

func devToolSource(ctx *typictx.Context, configuration ProjectConfiguration, out typiast.Outcome) *golang.SourceCode {
	return golang.NewSourceCode("main").
		AddImport(devToolSideEffects(ctx)...).
		AddStruct(configuration.Struct).
		AddConstructorFunction(configuration.Constructors...).
		AddConstructors(out.Autowires...).
		AddMockTargets(out.Automocks...).
		AddTestTargets(out.Packages...)
}

func appSource(ctx *typictx.Context, configuration ProjectConfiguration, out typiast.Outcome) *golang.SourceCode {
	return golang.NewSourceCode("main").
		AddImport(appSideEffects(ctx)...).
		AddStruct(configuration.Struct).
		AddConstructorFunction(configuration.Constructors...).
		AddConstructors(out.Autowires...).
		AddMockTargets(out.Automocks...).
		AddTestTargets(out.Packages...)
}
