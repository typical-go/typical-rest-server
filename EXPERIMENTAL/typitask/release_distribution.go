package typitask

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

// ReleaseDistribution to release binary distribution
func ReleaseDistribution(ctx typictx.ActionContext) error {
	RunTest(ctx)
	GenerateReadme(ctx)

	goos := []string{"linux", "darwin"}
	goarch := []string{"amd64"}
	mainPackage := typienv.AppMainPackage()

	for _, os1 := range goos {
		os.Setenv("GOOS", os1)
		for _, arch := range goarch {
			// TODO: using ldflags
			binaryName := fmt.Sprintf("%s/%s_%s_%s_%s",
				typienv.Release(),
				ctx.Typical.BinaryNameOrDefault(),
				ctx.Typical.Version,
				os1,
				arch)
			os.Setenv("GOARCH", arch)

			log.Infof("Create release for %s/%s: %s", os1, arch, binaryName)
			bash.GoBuild(binaryName, mainPackage)
		}
	}

	return nil
}
