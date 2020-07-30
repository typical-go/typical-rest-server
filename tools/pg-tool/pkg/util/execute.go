package util

import (
	"os"

	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
)

// Execute pg-tool
func Execute(c *typgo.Context, args []string) error {
	if _, err := os.Stat(bin); os.IsNotExist(err) {
		if err := c.Execute(&execkit.GoBuild{
			Output:      bin,
			MainPackage: src,
		}); err != nil {
			return err
		}
	}
	return c.Execute(&execkit.Command{
		Name:   bin,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})
}
