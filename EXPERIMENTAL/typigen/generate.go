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
	return runn.Execute(
		typienv.WriteEnvIfNotExist(ctx),
		appGenerated(ctx, configuration, out),
		devToolGeneratead(ctx, configuration, out),
	)
}

func devToolGeneratead(ctx *typictx.Context, configuration ProjectConfiguration, out typiast.Report) error {
	devTarget := typienv.TypicalDevToolMainPackage() + "/generated.go"
	sourceCode := golang.NewSourceCode("main").
		AddImport(devToolSideEffects(ctx)...).
		AddStruct(configuration.Struct).
		AddConstructorFunction(configuration.Constructors...).
		AddConstructors(out.Autowires...).
		AddMockTargets(out.Automocks...).
		AddTestTargets(out.Packages...)
	return runn.Execute(
		sourceCode.Cook(devTarget),
		bash.GoImports(devTarget),
	)
}

func appGenerated(ctx *typictx.Context, configuration ProjectConfiguration, report typiast.Report) error {
	appTarget := typienv.AppMainPackage() + "/generated.go"
	sourceCode := golang.NewSourceCode("main").
		AddImport(appSideEffects(ctx)...).
		AddStruct(configuration.Struct).
		AddConstructorFunction(configuration.Constructors...).
		AddConstructors(report.Autowires...).
		AddMockTargets(report.Automocks...).
		AddTestTargets(report.Packages...)

	return runn.Execute(
		sourceCode.Cook(appTarget),
		bash.GoImports(appTarget),
	)
}
