package typigen

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiparser"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/golang"
)

// GenDependecies generate dependecies
func GenDependecies(ctx *typictx.Context, srcCode *golang.SourceCode) (err error) {
	proj, err := typiparser.Parse("app")
	if err != nil {
		return
	}
	srcCode.AddConstructors(proj.Autowires...)
	srcCode.AddMockTargets(proj.Automocks...)
	srcCode.AddTestTargets(proj.Packages...)
	return
}
