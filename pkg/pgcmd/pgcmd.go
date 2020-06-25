package pgcmd

import (
	"os"

	"github.com/typical-go/typical-go/pkg/execkit"
)

// Param for command
type Param struct {
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

// PgTool is command for pg-tool
func PgTool(p *Param) *execkit.Command {
	return &execkit.Command{
		Name: p.Name,
		Args: []string{
			p.Action,
			"-host=" + p.Host,
			"-port=" + p.Port,
			"-user=" + p.User,
			"-password=" + p.Password,
			"-db-name=" + p.DBName,
			"-migration-src=" + p.MigrationSrc,
			"-seed-src=" + p.SeedSrc,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
}

// Console postgrs
func Console(p *Param) *execkit.Command {
	os.Setenv("PGPASSWORD", p.Password)

	// TODO: using `docker -it` for psql

	return &execkit.Command{
		Name: "psql",
		Args: []string{
			"-h", p.Host,
			"-p", p.Port,
			"-U", p.User,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
}
