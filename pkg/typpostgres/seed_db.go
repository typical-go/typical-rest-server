package typpostgres

import (
	"database/sql"
	"io/ioutil"

	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/urfave/cli/v2"
)

func cmdSeedDB(c *typbuildtool.Context) *cli.Command {
	return &cli.Command{
		Name:   "seed",
		Usage:  "Data seeding",
		Action: c.ActionFunc("PG", seedDB),
	}
}

func seedDB(c *typbuildtool.CliContext) (err error) {
	var (
		db  *sql.DB
		cfg *Config
	)

	if cfg, err = retrieveConfig(); err != nil {
		return
	}

	if db, err = sql.Open("postgres", dataSource(cfg)); err != nil {
		return
	}
	defer db.Close()

	files, _ := ioutil.ReadDir(DefaultSeedSource)
	for _, f := range files {
		sqlFile := DefaultSeedSource + "/" + f.Name()
		c.Infof("Execute seed '%s'", sqlFile)
		var b []byte
		if b, err = ioutil.ReadFile(sqlFile); err != nil {
			c.Warn(err.Error())
			continue
		}
		if _, err = db.ExecContext(c.Context, string(b)); err != nil {
			c.Warn(err.Error())
		}
	}
	return
}
