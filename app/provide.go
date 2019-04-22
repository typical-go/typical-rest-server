package app

import (
	"github.com/imantung/typical-go-server/app/cntrl"
	"github.com/imantung/typical-go-server/app/repo"
	"github.com/imantung/typical-go-server/config"
	"go.uber.org/dig"
)

func container() *dig.Container {
	container := dig.New()
	container.Provide(newServer)
	container.Provide(config.LoadConfig)
	container.Provide(connectDB)
	container.Provide(cntrl.NewBookController)
	container.Provide(repo.NewBookRepository)

	return container
}
