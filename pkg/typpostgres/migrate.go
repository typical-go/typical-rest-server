package typpostgres

import (
	"github.com/golang-migrate/migrate"
	log "github.com/sirupsen/logrus"
)

func (m *Module) migrate(cfg Config) (err error) {
	var (
		migration *migrate.Migrate
		sourceURL = "file://" + m.MigrationSource
	)
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, m.dataSource(cfg)); err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}
