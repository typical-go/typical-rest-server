package typpostgres

import (
	"github.com/golang-migrate/migrate"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

func (m *Module) rollbackCmd(c *typcore.BuildContext) *cli.Command {
	return &cli.Command{
		Name:   "rollback",
		Usage:  "Rollback Database",
		Action: c.ActionFunc(m.rollback),
	}
}

func (m *Module) rollback(cfg Config) (err error) {
	var migration *migrate.Migrate
	sourceURL := "file://" + migrationSrc
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, m.dataSource(cfg)); err != nil {
		return
	}
	defer migration.Close()
	return migration.Down()
}
