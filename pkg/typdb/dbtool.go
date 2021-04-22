package typdb

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/typical-go/typical-go/pkg/typgo"
)

type (
	DBTool struct {
		DBConn
		Name         string
		EnvKeys      *EnvKeys
		MigrationSrc string
		SeedSrc      string
		CreateFormat string
		DropFormat   string
	}
	DBConn interface {
		Connect(*Config) (*sql.DB, error)
		ConnectAdmin(*Config) (*sql.DB, error)
		Migrate(src string, cfg *Config) (*migrate.Migrate, error)
	}
)

var _ (typgo.Tasker) = (*DBTool)(nil)

// Task for postgres
func (t *DBTool) Task() *typgo.Task {
	return &typgo.Task{
		Name:  t.Name,
		Usage: fmt.Sprintf("%s database tool", t.Name),
		SubTasks: []*typgo.Task{
			{Name: "create", Usage: "Create database", Action: typgo.NewAction(t.CreateDB)},
			{Name: "drop", Usage: "Drop database", Action: typgo.NewAction(t.DropDB)},
			{Name: "migrate", Usage: "Migrate database", Action: typgo.NewAction(t.MigrateDB)},
			{Name: "migration", Usage: "Create Migration file", Action: typgo.NewAction(t.MigrationFile)},
			{Name: "rollback", Usage: "Rollback database", Action: typgo.NewAction(t.RollbackDB)},
			{Name: "seed", Usage: "Seed database", Action: typgo.NewAction(t.SeedDB)},
		},
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
