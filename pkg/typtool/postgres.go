package typtool

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/typical-go/typical-go/pkg/typgo"

	// load migration file
	_ "github.com/golang-migrate/migrate/source/file"
)

type (
	Postgres struct {
		Name         string
		EnvKeys      *DBEnvKeys
		MigrationSrc string
		SeedSrc      string
		DockerName   string
	}
)

var _ (typgo.Tasker) = (*Postgres)(nil)

// Task for postgres
func (t *Postgres) Task() *typgo.Task {
	return &typgo.Task{
		Name:  t.Name,
		Usage: "postgres database tool",
		SubTasks: []*typgo.Task{
			{Name: "create", Usage: "Create database", Action: typgo.NewAction(t.CreateDB)},
			{Name: "drop", Usage: "Drop database", Action: typgo.NewAction(t.DropDB)},
			{Name: "migrate", Usage: "Migrate database", Action: typgo.NewAction(t.MigrateDB)},
			{Name: "migration", Usage: "Create Migration file", Action: typgo.NewAction(t.MigrationFile)},
			{Name: "rollback", Usage: "Rollback database", Action: typgo.NewAction(t.RollbackDB)},
			{Name: "seed", Usage: "Seed database", Action: typgo.NewAction(t.SeedDB)},
			{Name: "console", Usage: "Postgres console", Action: typgo.NewAction(t.Console)},
		},
	}
}

// Console interactice for postgres
func (t *Postgres) Console(c *typgo.Context) error {
	cfg := t.EnvKeys.GetConfig()
	os.Setenv("PGPASSWORD", cfg.DBPass)
	return c.Execute(&typgo.Bash{
		Name: "docker",
		Args: []string{
			"exec", "-it", t.dockerName(),
			"psql",
			"-h", cfg.Host,
			"-p", cfg.Port,
			"-U", cfg.DBUser,
			"-d", cfg.DBName,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})
}

func (t *Postgres) dockerName() string {
	dockerName := t.DockerName
	if dockerName == "" {
		dockerName = typgo.ProjectName + "-" + t.Name
	}
	return dockerName
}

// CreateDB create database
func (t *Postgres) CreateDB(c *typgo.Context) error {
	cfg := t.EnvKeys.GetConfig()
	conn, err := openPostgresForAdmin(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf("CREATE DATABASE \"%s\"", cfg.DBName)
	c.Infof("%s: %s\n", t.Name, q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// DropDB delete database
func (t *Postgres) DropDB(c *typgo.Context) error {
	cfg := t.EnvKeys.GetConfig()
	conn, err := openPostgresForAdmin(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf("DROP DATABASE IF EXISTS \"%s\"", cfg.DBName)
	c.Infof("%s: %s\n", t.Name, q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// MigrateDB migrate database
func (t *Postgres) MigrateDB(c *typgo.Context) error {
	c.Infof("%s: Migrate '%s'\n", t.Name, t.MigrationSrc)
	migration, err := t.createMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

// RollbackDB rollback database
func (t *Postgres) RollbackDB(c *typgo.Context) error {
	c.Infof("%s: Rollback '%s'\n", t.Name, t.MigrationSrc)
	migration, err := t.createMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}

// SeedDB seed database
func (t *Postgres) SeedDB(c *typgo.Context) error {
	cfg := t.EnvKeys.GetConfig()
	db, err := openPostgres(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	files, _ := ioutil.ReadDir(t.SeedSrc)
	for _, f := range files {
		filename := fmt.Sprintf("%s/%s", t.SeedSrc, f.Name())
		c.Infof("%s: Seed '%s'\n", t.Name, filename)
		b, _ := ioutil.ReadFile(filename)
		_, err = db.ExecContext(c.Ctx(), string(b))
		if err != nil {
			return err
		}
	}
	return nil
}

// MigrationFile seed database
func (t *Postgres) MigrationFile(c *typgo.Context) error {
	args := c.Args().Slice()
	if len(args) < 1 {
		args = []string{"migration"}
	}
	for _, arg := range args {
		CreateMigrationFile(t.MigrationSrc, arg)
	}
	return nil
}

func (t *Postgres) createMigration() (*migrate.Migrate, error) {
	cfg := t.EnvKeys.GetConfig()
	db, err := openPostgres(cfg)
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	return migrate.NewWithDatabaseInstance("file://"+t.MigrationSrc, "postgres", driver)
}
