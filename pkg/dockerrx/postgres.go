package dockerrx

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typdocker"
)

// Postgres docker recipe
type Postgres struct {
	Version  string
	Name     string
	Image    string
	User     string
	Password string
	Port     int
}

var _ typdocker.Composer = (*Postgres)(nil)

// Compose to return the recipe
func (p *Postgres) Compose() (*typdocker.Recipe, error) {
	return &typdocker.Recipe{
		Version: p.Version,
		Services: typdocker.Services{
			p.Name: typdocker.Service{
				Image: p.Image,
				Environment: map[string]string{
					"POSTGRES":          p.User,
					"POSTGRES_PASSWORD": p.Password,
					"PGDATA":            "/data/postgres",
				},
				Volumes: []string{
					fmt.Sprintf("%s:/data/postgres", p.Name),
				},
				Ports: []string{
					fmt.Sprintf("%d:5432", p.Port),
				},
				Networks: []string{
					p.Name,
				},
				Restart: "unless-stopped",
			},
		},
		Networks: typdocker.Networks{
			p.Name: typdocker.Network{
				Driver: "bridge",
			},
		},
		Volumes: typdocker.Volumes{
			p.Name: nil,
		},
	}, nil
}
