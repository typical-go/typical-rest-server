package typical

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/urfave/cli/v2"

	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/pkg/dockerrx"
	"github.com/typical-go/typical-rest-server/pkg/typpg"
)

type (
	pgUtility struct{}
	pgDocker  struct {
		name string
	}
)

//
// pgUtility
//

var _ typgo.Utility = (*pgUtility)(nil)

func (*pgUtility) Commands(c *typgo.BuildCli) []*cli.Command {
	util := &typpg.Utility{
		Name:         "pg",
		MigrationSrc: "scripts/db/migration",
		SeedSrc:      "scripts/db/seed",
		ConfigName:   "PG",
	}
	return util.Commands(c)
}

//
// pgDocker
//

var _ typdocker.Composer = (*pgDocker)(nil)

func (p *pgDocker) Compose() (*typdocker.Recipe, error) {
	var cfg typpg.Config
	if err := typgo.ProcessConfig("PG", &cfg); err != nil {
		return nil, fmt.Errorf("pg-docker: " + err.Error())
	}

	pg := &dockerrx.Postgres{
		Version:  typdocker.V3,
		Name:     p.name,
		User:     cfg.User,
		Password: cfg.Password,
		Port:     cfg.Port,
	}
	return pg.Compose()
}
