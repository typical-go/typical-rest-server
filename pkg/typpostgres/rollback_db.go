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
		Action: rollbackDBAction(c),
	}
}

func rollbackDBAction(c *typbuildtool.Context) cli.ActionFunc {
	return func(cliCtx *cli.Context) (err error) {
		return rollbackDB(c.BuildContext(cliCtx))
	}
}

func rollbackDB(c *typbuildtool.BuildContext) (err error) {
	var (
		migration *migrate.Migrate
		cfg       *Config
	)

	if cfg, err = retrieveConfig(c); err != nil {
		return
	}

	sourceURL := "file://" + DefaultMigrationSource
	c.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, dataSource(cfg)); err != nil {
		return
	}
	defer migration.Close()
	return migration.Down()
}
