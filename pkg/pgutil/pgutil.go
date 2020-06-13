package pgutil

import (
	"context"
	"os"

	"github.com/typical-go/typical-go/pkg/execkit"
)

// PgUtil command
type PgUtil struct {
	Name         string
	Action       string
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	MigrationSrc string
	SeedSrc      string
}

// Command of pgUtil
func (c *PgUtil) Command() *execkit.Command {
	return &execkit.Command{
		Name: c.Name,
		Args: []string{
			c.Action,
			"-host=" + c.Host,
			"-port=" + c.Port,
			"-user=" + c.User,
			"-password=" + c.Password,
			"-db-name=" + c.DBName,
			"-migration-src=" + c.MigrationSrc,
			"-seed-src=" + c.SeedSrc,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
}

// Run command
func (c PgUtil) Run(ctx context.Context) error {
	return c.Command().Run(ctx)
}

func (c PgUtil) String() string {
	return c.Command().String()
}
