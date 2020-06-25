package dockerrx

import (
	"github.com/typical-go/typical-go/pkg/typdocker"
)

// Postgres docker recipe
type Postgres struct {
	Version  string
	Name     string
	Image    string
	User     string
	Password string
	Port     string
}

var _ typdocker.Composer = (*Postgres)(nil)

// Compose to return the recipe
func (p *Postgres) Compose() (*typdocker.Recipe, error) {
	if p.Version == "" {
		p.Version = typdocker.V3
	}
	if p.Name == "" {
		p.Name = "pg"
	}
	if p.Image == "" {
		p.Image = "postgres"
	}

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
					p.Name + ":/data/postgres",
				},
				Ports: []string{
					p.Port + ":5432",
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
