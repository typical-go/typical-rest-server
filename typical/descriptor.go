package typical

import (
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/typrails"
	"github.com/typical-go/typical-rest-server/pkg/typreadme"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
)

var (
	application = app.New()
	readme      = typreadme.New()
	rails       = typrails.New()
	server      = typserver.New()
	redis       = typredis.New()
	postgres    = typpostgres.New().WithDBName("sample")
	docker      = typdocker.New().WithComposers(postgres, redis)

	// Descriptor of Typical REST Server
	Descriptor = typcore.ProjectDescriptor{
		Name:        "Typical REST Server",
		Description: "Example of typical and scalable RESTful API Server for Go",
		Version:     "0.8.15",
		Package:     "github.com/typical-go/typical-rest-server",

		App: typcore.NewApp().
			WithEntryPoint(application).
			WithProvide(
				server,
				redis,
				postgres,
			).
			WithDestroy(
				server,
				redis,
				postgres,
			).
			WithPrepare(
				redis,
				postgres,
			),

		BuildCommands: []typcore.BuildCommander{
			docker,
			readme,
			postgres,
			redis,
			rails,
		},

		Configuration: typcore.NewConfiguration().
			WithConfigure(
				application,
				server,
				redis,
				postgres,
			),

		Releaser: typrls.New().WithPublisher(
			typrls.GithubPublisher("typical-go", "typical-rest-server"),
		),
	}
)
