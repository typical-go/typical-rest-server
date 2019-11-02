package typicli

import (
	"github.com/typical-go/typical-rest-server/pkg/utility/envkit"
	"github.com/urfave/cli"
)

// LoadEnvFile is cli version of LoadEnvFile
func LoadEnvFile(ctx *cli.Context) (err error) {
	return envkit.LoadEnvFile()
}
