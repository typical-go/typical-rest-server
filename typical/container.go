package typical

import (
	"github.com/typical-go/typical-rest-server/app/controller"
	"github.com/typical-go/typical-rest-server/app/entity"
	"github.com/typical-go/typical-rest-server/app/server"
	"github.com/typical-go/typical-rest-server/typical/infra"

	"go.uber.org/dig"
)

// Container dependency container
func Container() *dig.Container {
	container := dig.New()
	container.Provide(server.NewServer)

	// config
	container.Provide(LoadConfig)
	container.Provide(LoadPostgresConfig)

	// infra
	container.Provide(infra.ConnectPostgres)
	container.Provide(infra.CreatePostgresInfra)

	// controller
	container.Provide(controller.NewBookController)

	// entity
	container.Provide(entity.NewBookRepository)

	return container
}
