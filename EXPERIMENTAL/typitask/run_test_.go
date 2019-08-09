package typitask

import (
	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

// RunTest to run automate testing
func RunTest(ctx *typictx.ActionContext) error {
	return runn.Execute(
		bash.GoModTidy(),
		bash.GoTest(ctx.Typical.TestTargets),
	)
}
