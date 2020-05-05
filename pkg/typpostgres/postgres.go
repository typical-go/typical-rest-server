package typpostgres

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/urfave/cli/v2"
)

var (
	_ typcfg.Configurer = (*Postgres)(nil)

	// Default instance
	Default = Create(nil)
)

// Postgres module
type Postgres struct {
	typdocker.Composer
	typbuildtool.Utility
	*typapp.Module
	*typcfg.Configuration
}

// Create instance of psotgres module
func Create(s *Setting) *Postgres {
	if s == nil {
		s = &Setting{}
	}
	return &Postgres{
		Composer: recipe(s),

		Utility: typbuildtool.NewUtility(Commands),

		Module: typapp.NewModule().
			Provide(
				typapp.NewConstructor("", Connect),
			).
			Destroy(
				typapp.NewDestructor(Disconnect),
			).
			Prepare(
				typapp.NewPreparation(Ping),
			),

		Configuration: configuration(s),
	}
}

// Commands of postgres
func (p *Postgres) Commands(ctx *typbuildtool.Context) []*cli.Command {
	return p.Utility.Commands(ctx)
}

// Configurations of postgres
func (p *Postgres) Configurations() []*typcfg.Configuration {
	return p.Configuration.Configurations()
}

func configuration(s *Setting) *typcfg.Configuration {
	return &typcfg.Configuration{
		Name: ConfigName(s),
		Spec: &Config{
			DBName:   DBName(s),
			User:     User(s),
			Password: Password(s),
			Host:     Host(s),
			Port:     Port(s),
		},
	}
}

func recipe(s *Setting) *typdocker.Recipe {
	name := DockerName(s)
	image := DockerImage(s)

	return &typdocker.Recipe{
		Version: typdocker.V3,
		Services: typdocker.Services{
			name: typdocker.Service{
				Image: image,
				Environment: map[string]string{
					"POSTGRES":          User(s),
					"POSTGRES_PASSWORD": Password(s),
					"PGDATA":            "/data/postgres",
				},
				Volumes: []string{
					"postgres:/data/postgres",
				},
				Ports: []string{
					fmt.Sprintf("%d:5432", Port(s)),
				},
				Networks: []string{
					name,
				},
				Restart: "unless-stopped",
			},
		},
		Networks: typdocker.Networks{
			name: typdocker.Network{
				Driver: "bridge",
			},
		},
		Volumes: typdocker.Volumes{
			name: nil,
		},
	}
}

// Commands of module
func Commands(c *typbuildtool.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "postgres",
			Aliases: []string{"pg"},
			Usage:   "Postgres utility",
			Subcommands: []*cli.Command{
				cmdCreateDB(c),
				cmdDropDB(c),
				cmdMigrateDB(c),
				cmdRollbackDB(c),
				cmdSeedDB(c),
				cmdResetDB(c),
				cmdConsole(c),
			},
		},
	}
}
