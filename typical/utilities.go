package typical

import (
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/urfave/cli/v2"
)

type Utilities []typgo.Utility

// Commands of Utilities
func (u Utilities) Commands(b *typgo.BuildTool) (cmds []*cli.Command) {
	for _, utility := range u {
		cmds = append(cmds, utility.Commands(b)...)
	}
	return
}
