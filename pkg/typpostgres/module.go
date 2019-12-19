package typpostgres

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/utility/common"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/urfave/cli/v2"
)

const (
	migrationSrc = "scripts/db/migration"
	seedSrc      = "scripts/db/seed"
)

// Config is postgres configuration
type Config struct {
	DBName   string `required:"true"`
	User     string `required:"true" default:"postgres"`
	Password string `required:"true" default:"pgpass"`
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
}

// Module of postgres
type Module struct {
	DBName string
}

// BuildCommands of module
func (m Module) BuildCommands(c *typcore.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "postgres",
			Aliases: []string{"pg"},
			Usage:   "Postgres Database Tool",
			Before: func(ctx *cli.Context) error {
				return common.LoadEnvFile()
			},
			Subcommands: []*cli.Command{
				m.createCmd(c),
				m.dropCmd(c),
				m.migrateCmd(c),
				m.rollbackCmd(c),
				m.seedCmd(c),
				m.resetCmd(c),
				m.consoleCmd(c),
			},
		},
	}
}

// Provide the dependencies
func (m Module) Provide() []interface{} {
	return []interface{}{
		m.open,
	}
}

// Prepare the module
func (m Module) Prepare() []interface{} {
	return []interface{}{
		m.ping,
	}
}

// Configure the module
func (m Module) Configure() (prefix string, spec, loadFn interface{}) {
	prefix = "PG"
	spec = &Config{
		DBName: m.DBName,
	}
	loadFn = func(loader typcore.ConfigLoader) (cfg Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}

// Destroy dependencies
func (m Module) Destroy() []interface{} {
	return []interface{}{
		m.closeConnection,
	}
}

// DockerCompose template
func (m Module) DockerCompose() typdocker.Compose {
	return typdocker.Compose{
		Services: map[string]interface{}{
			"postgres": typdocker.Service{
				Image: "postgres",
				Environment: map[string]string{
					"POSTGRES":          "${PG_USER:-postgres}",
					"POSTGRES_PASSWORD": "${PG_PASSWORD:-pgpass}",
					"PGDATA":            "/data/postgres",
				},
				Volumes:  []string{"postgres:/data/postgres"},
				Ports:    []string{"${PG_PORT:-5432}:5432"},
				Networks: []string{"postgres"},
				Restart:  "unless-stopped",
			},
		},
		Networks: map[string]interface{}{
			"postgres": typdocker.Network{
				Driver: "bridge",
			},
		},
		Volumes: map[string]interface{}{
			"postgres": nil,
		},
	}
}

func (m Module) open(cfg Config) (*sql.DB, error) {
	return sql.Open("postgres", m.dataSource(cfg))
}

func (m Module) ping(db *sql.DB) error {
	log.Info("Ping to Postgres")
	return db.Ping()
}

func (Module) closeConnection(db *sql.DB) error {
	return db.Close()
}

func (Module) dataSource(c Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func (Module) adminDataSource(c Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, "template1")
}
