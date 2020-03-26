package typpostgres

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typcfg"
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

	// DefaultConfigName is default value of ConfigName
	DefaultConfigName = "PG"
)

// Module of postgres
type Module struct {
	dbName          string
	user            string
	password        string
	host            string
	port            int
	dockerImage     string
	dockerName      string
	migrationSource string
	seedSource      string
	configName      string
}

// New postgres module
func New() *Module {
	return &Module{
		user:            defaultUser,
		password:        defaultPassword,
		host:            defaultHost,
		port:            defaultPort,
		dockerImage:     defaultDockerImage,
		dockerName:      defaultDockerName,
		migrationSource: defaultMigrationSource,
		seedSource:      defaultSeedSource,
		configName:      DefaultConfigName,
	}
}

// WithDBName return module with new database name
func (m *Module) WithDBName(dbname string) *Module {
	m.dbName = dbname
	return m
}

// WithUser to return module with new user
func (m *Module) WithUser(user string) *Module {
	m.user = user
	return m
}

// WithHost return module with new host
func (m *Module) WithHost(host string) *Module {
	m.host = host
	return m
}

// WithPort return module with new port
func (m *Module) WithPort(port int) *Module {
	m.port = port
	return m
}

// WithPassword return module with new password
func (m *Module) WithPassword(password string) *Module {
	m.password = password
	return m
}

// WithDockerName to return module with new docker name
func (m *Module) WithDockerName(dockerName string) *Module {
	m.dockerName = dockerName
	return m
}

// WithDockerImage return module with new docker image
func (m *Module) WithDockerImage(dockerImage string) *Module {
	m.dockerImage = dockerImage
	return m
}

// WithMigrationSource return module with new migration source
func (m *Module) WithMigrationSource(migrationSource string) *Module {
	m.migrationSource = migrationSource
	return m
}

// WithSeedSource return module with new seed source
func (m *Module) WithSeedSource(seedSource string) *Module {
	m.seedSource = seedSource
	return m
}

// WithConfigName return module with new configName
func (m *Module) WithConfigName(configName string) *Module {
	m.configName = configName
	return m
}

// Configure the module
func (m *Module) Configure() *typcfg.Configuration {
	return typcfg.NewConfiguration(m.configName, &Config{
		DBName:   m.dbName,
		User:     m.user,
		Password: m.password,
		Host:     m.host,
		Port:     m.port,
	})
}

// Provide the dependencies
func (m *Module) Provide() []*typapp.Constructor {
	return []*typapp.Constructor{
		typapp.NewConstructor(Connect),
	}
}

// Prepare the module
func (m *Module) Prepare() []*typapp.Preparation {
	return []*typapp.Preparation{
		typapp.NewPreparation(Ping),
	}
}

// Destroy dependencies
func (m *Module) Destroy() []*typapp.Destruction {
	return []*typapp.Destruction{
		typapp.NewDestruction(Disconnect),
	}
}

// Connect to postgres server
func Connect(cfg *Config) (pgDB *DB, err error) {
	var db *sql.DB
	if db, err = sql.Open("postgres", dataSource(cfg)); err != nil {
		err = fmt.Errorf("Posgres: Connect: %w", err)
	}
	pgDB = NewDB(db)
	return
}

// Disconnect to postgres server
func Disconnect(db *DB) (err error) {
	if err = db.Close(); err != nil {
		return fmt.Errorf("Postgres: Disconnect: %w", err)
	}
	return
}

// Ping to postgres server
func Ping(db *DB) (err error) {
	if err = db.Ping(); err != nil {
		return fmt.Errorf("Postgres: Ping: %w", err)
	}
	return
}

func dataSource(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func adminDataSource(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, "template1")
}
