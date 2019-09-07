package typigen

import (
	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiparser"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/golang"
)

// Generate all
func Generate(ctx *typictx.Context) (err error) {
	proj, err := typiparser.Parse("app")
	if err != nil {
		return
	}
	configuration := configuration(ctx)
	appOut := typienv.AppMainPackage() + "/generated.go"
	devOut := typienv.TypicalDevToolMainPackage() + "/generated.go"
	return runn.Execute(
		typienv.WriteEnvIfNotExist(ctx),
		appSource(ctx, configuration, proj).Cook(appOut),
		devToolSource(ctx, configuration, proj).Cook(devOut),
		bash.GoImports(appOut),
		bash.GoImports(devOut),
	)
}

func devToolSource(ctx *typictx.Context, configuration ProjectConfiguration, proj typiparser.ProjectContext) *golang.SourceCode {
	return golang.NewSourceCode("main").
		AddImport(devToolSideEffects(ctx)...).
		AddStruct(configuration.Struct).
		AddConstructorFunction(configuration.Constructors...).
		AddConstructors(proj.Autowires...).
		AddMockTargets(proj.Automocks...).
		AddTestTargets(proj.Packages...)
}

func appSource(ctx *typictx.Context, configuration ProjectConfiguration, proj typiparser.ProjectContext) *golang.SourceCode {
	return golang.NewSourceCode("main").
		AddImport(appSideEffects(ctx)...).
		AddStruct(configuration.Struct).
		AddConstructorFunction(configuration.Constructors...).
		AddConstructors(proj.Autowires...).
		AddMockTargets(proj.Automocks...).
		AddTestTargets(proj.Packages...)
}
