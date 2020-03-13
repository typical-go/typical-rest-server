package typical

import (
	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-go/pkg/typreadme"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/typrails"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
	"github.com/typical-go/typical-rest-server/server"
)

// Modules that required for the project
var (
	serverApp = server.New().WithDebug(true)
	readme    = typreadme.New()
	rails     = typrails.New()
	redis     = typredis.New()
	postgres  = typpostgres.New().WithDBName("sample")

	docker = typdocker.New().WithComposers(
		postgres,
		redis,
	)
)
