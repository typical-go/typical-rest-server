package typpostgres

import (
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/typical-go/typical-go/pkg/typbuild"
	"github.com/typical-go/typical-go/pkg/typcfg"
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

// Module of postgres
type Module struct {
	DBName          string
	User            string
	Password        string
	Host            string
	Port            int
	DockerImage     string
	DockerName      string
	MigrationSource string
	SeedSource      string
}

// New postgres module
func New() *Module {
	return &Module{
		User:            defaultUser,
		Password:        defaultPassword,
		Host:            defaultHost,
		Port:            defaultPort,
		DockerImage:     defaultDockerImage,
		DockerName:      defaultDockerName,
		MigrationSource: defaultMigrationSource,
		SeedSource:      defaultSeedSource,
	}
}

// WithDBName to set database name
func (m *Module) WithDBName(dbname string) *Module {
	m.DBName = dbname
	return m
}

// WithUser to set user
func (m *Module) WithUser(user string) *Module {
	m.User = user
	return m
}

// WithHost to set host
func (m *Module) WithHost(host string) *Module {
	m.Host = host
	return m
}

// WithPort to set port
func (m *Module) WithPort(port int) *Module {
	m.Port = port
	return m
}

// WithPassword to set password
func (m *Module) WithPassword(password string) *Module {
	m.Password = password
	return m
}

// WithDockerName to set docker name
func (m *Module) WithDockerName(dockerName string) *Module {
	m.DockerName = dockerName
	return m
}

// WithDockerImage to set docker image
func (m *Module) WithDockerImage(dockerImage string) *Module {
	m.DockerImage = dockerImage
	return m
}

// WithMigrationSource to set migration source
func (m *Module) WithMigrationSource(migrationSource string) *Module {
	m.MigrationSource = migrationSource
	return m
}

// WithSeedSource to set seed source
func (m *Module) WithSeedSource(seedSource string) *Module {
	m.SeedSource = seedSource
	return m
}

// Configure the module
func (m *Module) Configure(loader typcfg.Loader) (prefix string, spec, constructor interface{}) {
	prefix = "PG"
	spec = &Config{
		DBName:   m.DBName,
		User:     m.User,
		Password: m.Password,
		Host:     m.Host,
		Port:     m.Port,
	}
	constructor = func() (cfg Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}

// BuildCommands of module
func (m *Module) BuildCommands(c *typbuild.Context) []*cli.Command {
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
func (m *Module) Provide() []interface{} {
	return []interface{}{
		m.connect,
	}
}

// Prepare the module
func (m *Module) Prepare() []interface{} {
	return []interface{}{
		m.ping,
	}
}

// Destroy dependencies
func (m *Module) Destroy() []interface{} {
	return []interface{}{
		m.disconnect,
	}
}
