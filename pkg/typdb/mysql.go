package typdb

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/typical-go/typical-go/pkg/typgo"

	// load migration file
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

type (
	MySQLTool struct {
		Name         string
		EnvKeys      *EnvKeys
		MigrationSrc string
		SeedSrc      string
		DockerName   string
	}
	MySQLHandler struct{}
)

//
// MySQL
//

var _ (typgo.Tasker) = (*MySQLTool)(nil)

// Task for postgress
func (t *MySQLTool) Task() *typgo.Task {
	return t.DBTool().Task()
}

func (t *MySQLTool) DBTool() *DBTool {
	if t.Name == "" {
		t.Name = "mysql"
	}
	return &DBTool{
		DBToolHandler: &MySQLHandler{},
		Name:          t.Name,
		EnvKeys:       t.EnvKeys,
		MigrationSrc:  t.MigrationSrc,
		SeedSrc:       t.SeedSrc,
		CreateFormat:  "CREATE DATABASE `%s`",
		DropFormat:    "DROP DATABASE IF EXISTS `%s`",
		DockerName:    t.DockerName,
	}
}

//
// MySQLHandler
//

var _ DBToolHandler = (*MySQLHandler)(nil)

func (MySQLHandler) Connect(c *Config) (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?tls=false&multiStatements=true",
		c.DBUser, c.DBPass, c.Host, c.Port, c.DBName,
	))
}

func (MySQLHandler) ConnectAdmin(c *Config) (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf(
		"root:%s@tcp(%s:%s)/?tls=false&multiStatements=true",
		c.DBPass, c.Host, c.Port,
	))
}

func (m MySQLHandler) Migrate(src string, cfg *Config) (*migrate.Migrate, error) {
	db, err := m.Connect(cfg)
	if err != nil {
		return nil, err
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, err
	}
	return migrate.NewWithDatabaseInstance(src, "mysql", driver)
}

// Console interactice for postgres
func (m MySQLHandler) Console(d *DBTool, c *typgo.Context) error {
	cfg := d.EnvKeys.Config()
	return c.Execute(&typgo.Bash{
		Name: "docker",
		Args: []string{
			"exec", "-it", d.DockerName,
			"mysql",
			"-h", "localhost",
			"-P", "3306",
			"-u", cfg.DBUser,
			fmt.Sprintf("-p%s", cfg.DBPass), // password flag can't be spaced
			cfg.DBName,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})
}
