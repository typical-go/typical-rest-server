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

// Postgres of postgres
type Postgres struct {
	DBName          string
	User            string
	Password        string
	Host            string
	Port            int
	DockerImage     string
	DockerName      string
	MigrationSource string
	SeedSource      string
	prefix          string
}

// New postgres module
func New() *Postgres {
	return &Postgres{
		User:            defaultUser,
		Password:        defaultPassword,
		Host:            defaultHost,
		Port:            defaultPort,
		DockerImage:     defaultDockerImage,
		DockerName:      defaultDockerName,
		MigrationSource: defaultMigrationSource,
		SeedSource:      defaultSeedSource,
		prefix:          "PG",
	}
}

// WithDBName to return module with new database name
func (m *Postgres) WithDBName(dbname string) *Postgres {
	m.DBName = dbname
	return m
}

// WithUser to return module with new user
func (m *Postgres) WithUser(user string) *Postgres {
	m.User = user
	return m
}

// WithHost to return module with new host
func (m *Postgres) WithHost(host string) *Postgres {
	m.Host = host
	return m
}

// WithPort to return module with new port
func (m *Postgres) WithPort(port int) *Postgres {
	m.Port = port
	return m
}

// WithPassword to return module with new password
func (m *Postgres) WithPassword(password string) *Postgres {
	m.Password = password
	return m
}

// WithDockerName to return module with new docker name
func (m *Postgres) WithDockerName(dockerName string) *Postgres {
	m.DockerName = dockerName
	return m
}

// WithDockerImage to return module with new docker image
func (m *Postgres) WithDockerImage(dockerImage string) *Postgres {
	m.DockerImage = dockerImage
	return m
}

// WithMigrationSource to return module with new migration source
func (m *Postgres) WithMigrationSource(migrationSource string) *Postgres {
	m.MigrationSource = migrationSource
	return m
}

// WithSeedSource to return module with new migration source
func (m *Postgres) WithSeedSource(seedSource string) *Postgres {
	m.SeedSource = seedSource
	return m
}

// Configure the module
func (m *Postgres) Configure(loader typcfg.Loader) *typcfg.Detail {
	return &typcfg.Detail{
		Prefix: m.prefix,
		Spec: &Config{
			DBName:   m.DBName,
			User:     m.User,
			Password: m.Password,
			Host:     m.Host,
			Port:     m.Port,
		},
		Constructor: func() (cfg Config, err error) {
			err = loader.Load(m.prefix, &cfg)
			return
		},
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
func (m *Postgres) Provide() []interface{} {
	return []interface{}{
		m.connect,
	}
}

// Prepare the module
func (m *Postgres) Prepare() []interface{} {
	return []interface{}{
		m.ping,
	}
}

// Destroy dependencies
func (m *Postgres) Destroy() []interface{} {
	return []interface{}{
		m.disconnect,
	}
}
