package typpostgres

import (
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/typical-go/typical-go/pkg/typbuild"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typdep"
	"github.com/urfave/cli/v2"
)

const (
	defaultUser            = "postgres"
	defaultPassword        = "pgpass"
	defaultHost            = "localhost"
	defaultPort            = 5432
	defaultDockerImage     = "postgres"
	defaultDockerName      = "postgres"
	defaultMigrationSource = "scripts/db/migration"
	defaultSeedSource      = "scripts/db/seed"
)

// Postgres of postgres
type Postgres struct {
	dbName          string
	user            string
	password        string
	host            string
	port            int
	dockerImage     string
	dockerName      string
	migrationSource string
	seedSource      string
	prefix          string
}

// New postgres module
func New() *Postgres {
	return &Postgres{
		user:            defaultUser,
		password:        defaultPassword,
		host:            defaultHost,
		port:            defaultPort,
		dockerImage:     defaultDockerImage,
		dockerName:      defaultDockerName,
		migrationSource: defaultMigrationSource,
		seedSource:      defaultSeedSource,
		prefix:          "PG",
	}
}

// WithDBName to return module with new database name
func (m *Postgres) WithDBName(dbname string) *Postgres {
	m.dbName = dbname
	return m
}

// WithUser to return module with new user
func (m *Postgres) WithUser(user string) *Postgres {
	m.user = user
	return m
}

// WithHost to return module with new host
func (m *Postgres) WithHost(host string) *Postgres {
	m.host = host
	return m
}

// WithPort to return module with new port
func (m *Postgres) WithPort(port int) *Postgres {
	m.port = port
	return m
}

// WithPassword to return module with new password
func (m *Postgres) WithPassword(password string) *Postgres {
	m.password = password
	return m
}

// WithDockerName to return module with new docker name
func (m *Postgres) WithDockerName(dockerName string) *Postgres {
	m.dockerName = dockerName
	return m
}

// WithDockerImage to return module with new docker image
func (m *Postgres) WithDockerImage(dockerImage string) *Postgres {
	m.dockerImage = dockerImage
	return m
}

// WithMigrationSource to return module with new migration source
func (m *Postgres) WithMigrationSource(migrationSource string) *Postgres {
	m.migrationSource = migrationSource
	return m
}

// WithSeedSource to return module with new migration source
func (m *Postgres) WithSeedSource(seedSource string) *Postgres {
	m.seedSource = seedSource
	return m
}

// Configure the module
func (m *Postgres) Configure(loader typcfg.Loader) *typcfg.Detail {
	return &typcfg.Detail{
		Prefix: m.prefix,
		Spec: &Config{
			DBName:   m.dbName,
			User:     m.user,
			Password: m.password,
			Host:     m.host,
			Port:     m.port,
		},
		Constructor: typdep.NewConstructor(
			func() (cfg Config, err error) {
				err = loader.Load(m.prefix, &cfg)
				return
			}),
	}
}

// BuildCommands of module
func (m *Postgres) BuildCommands(c *typbuild.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "postgres",
			Aliases: []string{"pg"},
			Usage:   "Postgres Database Tool",
			Before: func(ctx *cli.Context) error {
				return typcfg.LoadEnvFile()
			},
			Subcommands: []*cli.Command{
				{
					Name:   "create",
					Usage:  "Create New Database",
					Action: c.ActionFunc(m.create),
				},
				{
					Name:   "drop",
					Usage:  "Drop Database",
					Action: c.ActionFunc(m.drop),
				},
				{
					Name:   "migrate",
					Usage:  "Migrate Database",
					Action: c.ActionFunc(m.migrate),
				},
				{
					Name:   "rollback",
					Usage:  "Rollback Database",
					Action: c.ActionFunc(m.rollback),
				},
				{
					Name:   "seed",
					Usage:  "Data seeding",
					Action: c.ActionFunc(m.seed),
				},
				{
					Name:   "reset",
					Usage:  "Reset Database",
					Action: c.ActionFunc(m.reset),
				},
				{
					Name:   "console",
					Usage:  "PostgreSQL Interactive",
					Action: c.ActionFunc(m.console),
				},
			},
		},
	}
}

// Provide the dependencies
func (m *Postgres) Provide() []*typdep.Constructor {
	return []*typdep.Constructor{
		typdep.NewConstructor(m.connect),
	}
}

// Prepare the module
func (m *Postgres) Prepare() []*typdep.Invocation {
	return []*typdep.Invocation{
		typdep.NewInvocation(m.ping),
	}
}

// Destroy dependencies
func (m *Postgres) Destroy() []*typdep.Invocation {
	return []*typdep.Invocation{
		typdep.NewInvocation(m.disconnect),
	}
}
