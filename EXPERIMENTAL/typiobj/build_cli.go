package typiobj

import (
	"github.com/urfave/cli"
)

// BuildCLI responsible to give command
type BuildCLI interface {
	Command() cli.Command
}

// IsBuildCLI return true if object implementation of CommandLiner
func IsBuildCLI(obj interface{}) (ok bool) {
	_, ok = obj.(BuildCLI)
	return
}
