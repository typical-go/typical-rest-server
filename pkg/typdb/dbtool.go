package typdb

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/iancoleman/strcase"
	"github.com/typical-go/typical-go/pkg/typgo"
)

type (
	DBTool struct {
		DBToolHandler
		Name         string
		EnvKeys      *EnvKeys
		MigrationSrc string
		SeedSrc      string
		CreateFormat string
		DropFormat   string
		DockerName   string
	}
	DBToolHandler interface {
		Connect(*Config) (*sql.DB, error)
		ConnectAdmin(*Config) (*sql.DB, error)
		Migrate(src string, cfg *Config) (*migrate.Migrate, error)
		Console(*DBTool, *typgo.Context) error
	}
)

var _ (typgo.Tasker) = (*DBTool)(nil)

// Task for postgres
func (t *DBTool) Task() *typgo.Task {
	t.initDefault()
	task := &typgo.Task{
		Name:  t.Name,
		Usage: fmt.Sprintf("%s database tool", t.Name),
		SubTasks: []*typgo.Task{
			{Name: "create", Usage: "Create database", Action: typgo.NewAction(t.CreateDB)},
			{Name: "drop", Usage: "Drop database", Action: typgo.NewAction(t.DropDB)},
			{Name: "migrate", Usage: "Migrate database", Action: typgo.NewAction(t.MigrateDB)},
			{Name: "migration", Usage: "Create Migration file", Action: typgo.NewAction(t.MigrationFile)},
			{Name: "rollback", Usage: "Rollback database", Action: typgo.NewAction(t.RollbackDB)},
			{Name: "seed", Usage: "Seed database", Action: typgo.NewAction(t.SeedDB)},
			{Name: "console", Usage: "Database client console", Action: typgo.NewAction(t.Console)},
		},
	}
	return task
}

func (t *DBTool) initDefault() {
	if t.Name == "" {
		t.Name = "db"
	}
	if t.EnvKeys == nil {
		t.EnvKeys = EnvKeysWithPrefix(strcase.ToScreamingSnake(t.Name))
	}
	if t.MigrationSrc == "" {
		t.MigrationSrc = fmt.Sprintf("database/%s/migration", t.Name)
	}
	if t.SeedSrc == "" {
		t.SeedSrc = fmt.Sprintf("database/%s/seed", t.Name)
	}
	if t.DockerName == "" {
		t.DockerName = fmt.Sprintf("%s-%s", typgo.ProjectName, t.Name)
	}
}

// CreateDB create database
func (t *DBTool) CreateDB(c *typgo.Context) error {
	cfg := t.EnvKeys.Config()
	conn, err := t.ConnectAdmin(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf(t.CreateFormat, cfg.DBName)
	c.Infof("%s: %s\n", t.Name, q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// DropDB delete database
func (t *DBTool) DropDB(c *typgo.Context) error {
	cfg := t.EnvKeys.Config()
	conn, err := t.ConnectAdmin(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf(t.DropFormat, cfg.DBName)
	c.Infof("%s: %s\n", t.Name, q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// MigrateDB migrate database
func (t *DBTool) MigrateDB(c *typgo.Context) error {
	c.Infof("%s: Migrate '%s'\n", t.Name, t.MigrationSrc)
	m, err := t.Migrate("file://"+t.MigrationSrc, t.EnvKeys.Config())
	if err != nil {
		return err
	}
	defer m.Close()
	return m.Up()
}

// RollbackDB rollback database
func (t *DBTool) RollbackDB(c *typgo.Context) error {
	c.Infof("%s: Rollback '%s'\n", t.Name, t.MigrationSrc)
	m, err := t.Migrate("file://"+t.MigrationSrc, t.EnvKeys.Config())
	if err != nil {
		return err
	}
	defer m.Close()
	return m.Down()
}

// SeedDB seed database
func (t *DBTool) SeedDB(c *typgo.Context) error {
	cfg := t.EnvKeys.Config()
	db, err := t.Connect(cfg)
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
func (t *DBTool) MigrationFile(c *typgo.Context) error {
	args := c.Args().Slice()
	if len(args) < 1 {
		args = []string{"migration"}
	}
	for _, arg := range args {
		createMigration(t.MigrationSrc, arg)
	}
	return nil
}

func createMigration(migrationSrc, name string) {
	epoch := time.Now().Unix()
	upScript := fmt.Sprintf("%s/%d_%s.up.sql", migrationSrc, epoch, name)
	downScript := fmt.Sprintf("%s/%d_%s.down.sql", migrationSrc, epoch, name)

	if _, err := os.Create(upScript); err == nil {
		fmt.Println(upScript)
	}
	if _, err := os.Create(downScript); err == nil {
		fmt.Println(downScript)
	}
}

// Console interactice for postgres
func (t *DBTool) Console(c *typgo.Context) error {
	if t.DBToolHandler == nil {
		return nil
	}
	return t.DBToolHandler.Console(t, c)
}
