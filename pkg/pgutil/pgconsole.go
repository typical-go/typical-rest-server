package pgutil

import (
	"context"
	"os"

	"github.com/typical-go/typical-go/pkg/execkit"
)

// PgConsole to postgres
type PgConsole struct {
	Host     string
	Port     string
	User     string
	Password string
}

// Command to run postgres console
func (c *PgConsole) Command() *execkit.Command {
	os.Setenv("PGPASSWORD", c.Password)

	// TODO: using `docker -it` for psql

	return &execkit.Command{
		Name: "psql",
		Args: []string{
			"-h", c.Host,
			"-p", c.Port,
			"-U", c.User,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
}

// Run postgres console
func (c *PgConsole) Run(ctx context.Context) error {
	return c.Command().Run(ctx)
}

func (c PgConsole) String() string {
	return c.Command().String()
}
