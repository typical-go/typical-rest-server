package typpostgres

import (
	"github.com/golang-migrate/migrate"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/urfave/cli/v2"
)

func cmdRollbackDB(c *typbuildtool.Context) *cli.Command {
	return &cli.Command{
		Name:   "rollback",
		Usage:  "Rollback Database",
		Action: c.ActionFunc("PG", rollbackDB),
	}
}

func rollbackDB(c *typbuildtool.CliContext) (err error) {
	var (
		migration *migrate.Migrate
		cfg       *Config
	)

	if cfg, err = retrieveConfig(defaultConfigName); err != nil {
		return
	}

	sourceURL := "file://" + defaultMigrationSrc
	c.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, dataSource(cfg)); err != nil {
		return
	}
	defer migration.Close()
	return migration.Down()
}
