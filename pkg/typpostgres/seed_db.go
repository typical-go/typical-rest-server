package typpostgres

import (
	"database/sql"
	"io/ioutil"

	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/urfave/cli/v2"
)

func cmdSeedDB(c *typbuildtool.Context) *cli.Command {
	return &cli.Command{
		Name:  "seed",
		Usage: "Data seeding",
		Action: func(cliCtx *cli.Context) (err error) {
			return seedDB(c.BuildContext(cliCtx))
		},
	}
}

func seedDB(c *typbuildtool.BuildContext) (err error) {
	var conn *sql.DB
	var cfg *Config
	if cfg, err = retrieveConfig(c); err != nil {
		return
	}
	if conn, err = sql.Open("postgres", dataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	files, _ := ioutil.ReadDir(DefaultSeedSource)
	for _, f := range files {
		sqlFile := DefaultSeedSource + "/" + f.Name()
		c.Infof("Execute seed '%s'", sqlFile)
		var b []byte
		if b, err = ioutil.ReadFile(sqlFile); err != nil {
			c.Warn(err.Error())
			continue
		}
		if _, err = conn.Exec(string(b)); err != nil {
			c.Warn(err.Error())
		}
	}
	return
}
