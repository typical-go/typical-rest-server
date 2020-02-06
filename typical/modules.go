package typical

import (
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/typrails"
	"github.com/typical-go/typical-rest-server/pkg/typreadme"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
)

// Modules that required for the project
var (
	application = app.New()
	readme      = typreadme.New()
	rails       = typrails.New()
	server      = typserver.New().WithDebug(true)
	redis       = typredis.New()
	postgres    = typpostgres.New().WithDBName("sample")

	docker = typdocker.New().WithComposers(
		postgres,
		redis,
	)
)
