package app

import (
	"github.com/imantung/typical-go-server/app/controller"
	"github.com/imantung/typical-go-server/app/repository"
	"github.com/imantung/typical-go-server/config"
	"github.com/imantung/typical-go-server/db"
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

func container() *dig.Container {
	container := dig.New()
	container.Provide(newServer)
	container.Provide(config.LoadConfig)
	container.Provide(db.Connect)
	container.Provide(controller.NewBookController)
	container.Provide(repository.NewBookRepository)
	return container
}

func triggerAction(invokeFunc interface{}) interface{} {
	return func(ctx *cli.Context) error {
		container := container()
		container.Provide(ctx.Args)
		return container.Invoke(invokeFunc)
	}
}
