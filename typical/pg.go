package typical

import (
	"github.com/typical-go/typical-go/pkg/typgo"

	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/pkg/dockerrx"
	"github.com/typical-go/typical-rest-server/pkg/typpg"
)

type pgDocker struct{}

//
// pgDocker
//

var _ typdocker.Composer = (*pgDocker)(nil)

func (*pgDocker) DockerCompose() *typdocker.Recipe {
	var cfg typpg.Config
	if err := typgo.ProcessConfig("PG", &cfg); err != nil {
		panic("pg-docker: " + err.Error())
	}
	pg := &dockerrx.Postgres{
		Version:  typdocker.V3,
		Name:     "pg",
		Image:    "postgres",
		User:     cfg.User,
		Password: cfg.Password,
		Port:     cfg.Port,
	}
	return pg.DockerCompose()
}
