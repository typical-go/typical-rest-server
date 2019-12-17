package typpostgres

import (
	"github.com/golang-migrate/migrate"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

func (m Module) migrateCmd(c typcore.Cli) *cli.Command {
	return &cli.Command{
		Name:   "migrate",
		Usage:  "Migrate Database",
		Action: c.Action(m.migrate),
	}
}

func (m Module) migrate(cfg Config) (err error) {
	var migration *migrate.Migrate
	sourceURL := "file://" + migrationSrc
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, m.dataSource(cfg)); err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}
