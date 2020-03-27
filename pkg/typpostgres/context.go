package typpostgres

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/urfave/cli/v2"
)

// Context of postgress
type Context struct {
	*typbuildtool.Context
	*Module
	Cli *cli.Context
}

// Config of postgres
func (c *Context) Config() (cfg *Config, err error) {
	var v interface{}
	var ok bool

	if v, err = c.RetrieveConfig(ConfigName); err != nil {
		return
	}

	if cfg, ok = v.(*Config); !ok {
		return nil, fmt.Errorf("Postgres: Get config for '%s' but invalid type", ConfigName)
	}

	return
}
