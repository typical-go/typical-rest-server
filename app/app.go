package app

import (
	"github.com/imantung/typical-go-server/config"
	"github.com/imantung/typical-go-server/db"
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

var (
	container *dig.Container
)

// called once when package imported as go ecosystem
func init() {
	container = dig.New()
	container.Provide(config.NewConfig)
	container.Provide(db.Connect)
	container.Provide(newServer)
}

func triggerAction(function interface{}) interface{} {
	return func(ctx *cli.Context) error {
		return container.Invoke(function)
	}
}
