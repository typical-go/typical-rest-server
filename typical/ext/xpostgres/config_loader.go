package xpostgres

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/typical/appctx"
)

type ConfigLoader struct {
	appctx.ConfigDetail
}

func (l ConfigLoader) LoadFunc() interface{} {
	return func() (config Config, err error) {
		err = envconfig.Process(l.ConfigPrefix(), &config)
		return
	}
}
