package typicmd

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirelease"
)

// ReleaseDistribution to release binary distribution
func ReleaseDistribution(action *typictx.ActionContext) (err error) {
	if !action.Cli.Bool("no-test") {
		err = RunTest(action)
		if err != nil {
			return
		}
	}
	if !action.Cli.Bool("no-readme") {
		err = GenerateReadme(action)
		if err != nil {
			return
		}
	}
	return typirelease.Release(action.Context)
}
