package typicmd

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirelease"
)

func releaseDistribution(action *typictx.ActionContext) (err error) {
	if !action.Cli.Bool("no-test") {
		err = runTesting(action)
		if err != nil {
			return
		}
	}
	if !action.Cli.Bool("no-readme") {
		err = generateReadme(action)
		if err != nil {
			return
		}
	}
	binaries, changeLogs, err := typirelease.ReleaseDistribution(action.Release, action.Cli.Bool("force"))
	if err != nil {
		return
	}

	if !action.Cli.Bool("no-github") {
		typirelease.GithubRelease(binaries, changeLogs, action.Release)
	}
	return
}
