package typreadme

import (
	"github.com/typical-go/typical-go/pkg/typcore"
)

// ReadmeObject represent readme documentation
type ReadmeObject struct {
	Title               string
	Description         string
	ApplicationCommands CommandDetails
	OtherBuildCommands  CommandDetails
	Configs             typcore.ConfigDetails
}
