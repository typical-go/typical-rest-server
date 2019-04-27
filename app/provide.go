package app

import (
	"github.com/typical-go/typical-rest-server/app/cntrl"
	"github.com/typical-go/typical-rest-server/app/repo"
	"github.com/typical-go/typical-rest-server/config"
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
