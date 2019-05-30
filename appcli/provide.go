package appcli

import (
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/app/controller"
	"github.com/typical-go/typical-rest-server/app/entity"
	"github.com/typical-go/typical-rest-server/config"
	"go.uber.org/dig"
)

func container() *dig.Container {
	container := dig.New()
	container.Provide(app.NewServer)
	container.Provide(config.LoadConfig)
	container.Provide(connectDB)
	container.Provide(controller.NewBookController)
	container.Provide(entity.NewBookRepository)

	return container
}
