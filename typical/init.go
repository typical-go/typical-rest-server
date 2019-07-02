package typical

import (
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/app/controller"
	"github.com/typical-go/typical-rest-server/app/repository"
	"github.com/typical-go/typical-rest-server/typical/appctx"
	"github.com/typical-go/typical-rest-server/typical/ext/xpostgres"
)

// Context instance of Context
var Context *appctx.Context

func init() {
	Context = &appctx.Context{
		Name:           "Typical-RESTful-Server",
		Version:        "0.1.0",
		Description:    "Example of typical and scalable RESTful API Server for Go",
		ReadmeTemplate: readmeTemplate,

		ConfigLoader: ConfigLoader{
			ConfigDetail: appctx.NewConfigDetail("APP", &app.Config{}),
		},

		Constructors: []interface{}{
			app.NewServer,
			controller.NewBookController,
			repository.NewBookRepository,
		},

		Modules: []appctx.Module{
			xpostgres.NewModule(),
		},
	}

}
