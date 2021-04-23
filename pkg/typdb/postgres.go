package typdb

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/typical-go/typical-go/pkg/typgo"

	// load migration file
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

type (
	PostgresTool struct {
		Name         string
		EnvKeys      *EnvKeys
		MigrationSrc string
		SeedSrc      string
		DockerName   string
	}
	PGConn struct{}
)

//
// Postgres
//

var _ (typgo.Tasker) = (*PostgresTool)(nil)

// Task for postgres
func (t *PostgresTool) Task() *typgo.Task {
	dbtool := t.DBTool()
	subtasks := []*typgo.Task{
		{Name: "console", Usage: "Postgres console", Action: typgo.NewAction(t.Console)},
	}

	task := dbtool.Task()
	task.SubTasks = append(task.SubTasks, subtasks...)
	return task
}

func (t *PostgresTool) DBTool() *DBTool {
	return &DBTool{
		DBConn:       &PGConn{},
		Name:         t.Name,
		EnvKeys:      t.EnvKeys,
		MigrationSrc: t.MigrationSrc,
		SeedSrc:      t.SeedSrc,
		CreateFormat: "CREATE DATABASE \"%s\"",
		DropFormat:   "DROP DATABASE IF EXISTS \"%s\"",
	}
}

// Console interactice for postgres
func (t *PostgresTool) Console(c *typgo.Context) error {
	cfg := t.EnvKeys.Config()
	os.Setenv("PGPASSWORD", cfg.DBPass)
	return c.Execute(&typgo.Bash{
		Name: "docker",
		Args: []string{
			"exec", "-it", t.dockerName(),
			"psql",
			"-h", "localhost",
			"-p", "5432",
			"-U", cfg.DBUser,
			"-d", cfg.DBName,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})
}

func (t *PostgresTool) dockerName() string {
	dockerName := t.DockerName
	if dockerName == "" {
		dockerName = typgo.ProjectName + "-" + t.Name
	}
	return dockerName
}

//
// PGConn
//

var _ DBConn = (*PGConn)(nil)

func (PGConn) Connect(c *Config) (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPass, c.Host, c.Port, c.DBName,
	))
}

func (PGConn) ConnectAdmin(c *Config) (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/template1?sslmode=disable",
		c.DBUser, c.DBPass, c.Host, c.Port))
}

func (p PGConn) Migrate(src string, cfg *Config) (*migrate.Migrate, error) {
	db, err := p.Connect(cfg)
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	return migrate.NewWithDatabaseInstance(src, "postgres", driver)
}
