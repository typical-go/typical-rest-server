package app

import (
	"github.com/imantung/typical-go-server/app/controller"
	"github.com/imantung/typical-go-server/app/repository"
	"github.com/imantung/typical-go-server/config"
	"go.uber.org/dig"
)

func container() *dig.Container {
	container := dig.New()
	container.Provide(newServer)
	container.Provide(config.LoadConfig)
	container.Provide(connectDB)
	container.Provide(controller.NewBookController)
	container.Provide(repository.NewBookRepository)

	return container
}
