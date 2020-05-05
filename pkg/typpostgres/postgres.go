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
	// Default instance
	Default = Create()

	// DefaultDockerImage is default docker image for postgres
	DefaultDockerImage = "postgres"

	// DefaultDockerName is default docker name for postgres
	DefaultDockerName = "postgres"

	// DefaultMigrationSource is default migration source for postgres
	DefaultMigrationSource = "scripts/db/migration"

	// DefaultSeedSource is default seed source for postgres
	DefaultSeedSource = "scripts/db/seed"
)

var _ typcfg.Configurer = (*Postgres)(nil)

// Postgres module
type Postgres struct {
	typdocker.Composer
	typbuildtool.Utility
	*typapp.Module
	*typcfg.Configuration
}

// Create instance of psotgres module
func Create() *Postgres {
	return &Postgres{
		Composer: recipe(),

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

		Configuration: &typcfg.Configuration{
			Name: DefaultConfigName,
			Spec: DefaultConfig,
		},
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

func recipe() *typdocker.Recipe {
	return &typdocker.Recipe{
		Version: typdocker.V3,
		Services: typdocker.Services{
			DefaultDockerName: typdocker.Service{
				Image: DefaultDockerImage,
				Environment: map[string]string{
					"POSTGRES":          DefaultUser,
					"POSTGRES_PASSWORD": DefaultPassword,
					"PGDATA":            "/data/postgres",
				},
				Volumes:  []string{"postgres:/data/postgres"},
				Ports:    []string{fmt.Sprintf("%d:5432", DefaultPort)},
				Networks: []string{DefaultDockerName},
				Restart:  "unless-stopped",
			},
		},
		Networks: typdocker.Networks{
			DefaultDockerName: typdocker.Network{
				Driver: "bridge",
			},
		},
		Volumes: typdocker.Volumes{
			DefaultDockerName: nil,
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
