package typical

import (
	"os"

	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/pkg/dockerrx"
)

type pgDocker struct {
	name string
}

var _ typdocker.Composer = (*pgDocker)(nil)

func (p *pgDocker) Compose() (*typdocker.Recipe, error) {
	pg := &dockerrx.Postgres{
		Version:  typdocker.V3,
		Name:     p.name,
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		Port:     os.Getenv("PG_PORT"),
	}
	return pg.Compose()
}
