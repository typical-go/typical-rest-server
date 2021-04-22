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
	// MySQL for postgres
	MySQL struct {
		Name         string
		EnvKeys      *EnvKeys
		MigrationSrc string
		SeedSrc      string
		DockerName   string
	}
	PGConn struct{}
)

//
// MySQL
//

var _ (typgo.Tasker) = (*MySQL)(nil)

// Task for postgress
func (t *MySQL) Task() *typgo.Task {
	dbtool := t.DBTool()
	subtasks := []*typgo.Task{
		{Name: "console", Usage: "Postgres console", Action: typgo.NewAction(t.Console)},
	}

	task := dbtool.Task()
	task.SubTasks = append(task.SubTasks, subtasks...)
	return task
}

func (t *MySQL) DBTool() *DBTool {
	return &DBTool{
		DBConn:       &MySQLConn{},
		Name:         t.Name,
		EnvKeys:      t.EnvKeys,
		MigrationSrc: t.MigrationSrc,
		SeedSrc:      t.SeedSrc,
		CreateFormat: "CREATE DATABASE `%s`",
		DropFormat:   "DROP DATABASE IF EXISTS `%s`",
	}
}

// Console interactice for postgres
func (t *MySQL) Console(c *typgo.Context) error {
	cfg := t.EnvKeys.Config()
	return c.Execute(&typgo.Bash{
		Name: "docker",
		Args: []string{
			"exec", "-it", t.dockerName(),
			"mysql",
			"-h", "locahost",
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

func (t *MySQL) dockerName() string {
	dockerName := t.DockerName
	if dockerName == "" {
		dockerName = typgo.ProjectName + "-" + t.Name
	}
	return dockerName
}

//
// MySQLConn
//

var _ DBConn = (*MySQLConn)(nil)

func (MySQLConn) Connect(c *Config) (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?tls=false&multiStatements=true",
		c.DBUser, c.DBPass, c.Host, c.Port, c.DBName,
	))
}

func (MySQLConn) ConnectAdmin(c *Config) (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf(
		"root:%s@tcp(%s:%s)/?tls=false&multiStatements=true",
		c.DBPass, c.Host, c.Port,
	))
}

func (m MySQLConn) Migrate(src string, cfg *Config) (*migrate.Migrate, error) {
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
