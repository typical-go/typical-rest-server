package buildtool

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/buildtool/releaser"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

func releaseDistribution(action *typictx.ActionContext) (err error) {
	if !action.Cli.Bool("no-test") {
		err = runTesting(action)
		if err != nil {
			return
		}
	}
	binaries, changeLogs, err := releaser.ReleaseDistribution(action.Release, action.Cli.Bool("force"))
	if err != nil {
		return
	}

	if !action.Cli.Bool("no-github") {
		releaser.GithubRelease(binaries, changeLogs, action.Release)
	}
	return
}
