package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/urfave/cli/v2"
)

type (
	config struct {
		host         string
		port         string
		user         string
		password     string
		dbName       string
		migrationSrc string
		seedSrc      string
	}
)

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{Name: "host"},
		&cli.StringFlag{Name: "port"},
		&cli.StringFlag{Name: "user"},
		&cli.StringFlag{Name: "password"},
		&cli.StringFlag{Name: "db-name"},
		&cli.StringFlag{Name: "migration-src"},
		&cli.StringFlag{Name: "seed-src"},
	}
	app := &cli.App{
		Name: "MySQL Utility",
		Commands: []*cli.Command{
			{Name: "create", Flags: flags, Action: createDB},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func retrieveConfig(c *cli.Context) config {
	return config{
		host:         c.String("host"),
		port:         c.String("port"),
		user:         c.String("user"),
		password:     c.String("password"),
		dbName:       c.String("db-name"),
		migrationSrc: c.String("migration-src"),
		seedSrc:      c.String("seed-src"),
	}
}

func createDB(c *cli.Context) error {
	cfg := retrieveConfig(c)
	conn, err := sql.Open("mysql", connNoDB(cfg))
	if err != nil {
		return err
	}
	defer conn.Close()
	if err = conn.PingContext(c.Context); err != nil {
		return err
	}
	q := fmt.Sprintf("CREATE DATABASE `%s`", cfg.dbName)
	fmt.Println(q)
	_, err = conn.ExecContext(c.Context, q)
	return err
}

func dropDB(c *cli.Context) error {
	cfg := retrieveConfig(c)
	conn, err := sql.Open("mysql", connNoDB(cfg))
	if err != nil {
		return err
	}
	defer conn.Close()
	q := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", cfg.dbName)
	fmt.Println(q)
	_, err = conn.ExecContext(c.Context, q)
	return err
}

func migrateDB(c *cli.Context) error {
	cfg := retrieveConfig(c)
	sourceURL := "file://" + cfg.migrationSrc
	fmt.Printf("Migrate '%s'\n", cfg.migrationSrc)
	migration, err := migrate.New(sourceURL, conn(cfg))
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

func rollbackDB(c *cli.Context) error {
	cfg := retrieveConfig(c)
	sourceURL := "file://" + cfg.migrationSrc

	migration, err := migrate.New(sourceURL, conn(cfg))
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}

func seedDB(c *cli.Context) error {
	cfg := retrieveConfig(c)
	db, err := sql.Open("mysql", conn(cfg))
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Printf("Seed '%s'\n", cfg.seedSrc)
	files, _ := ioutil.ReadDir(cfg.seedSrc)
	for _, f := range files {
		sqlFile := cfg.seedSrc + "/" + f.Name()
		var b []byte
		if b, err = ioutil.ReadFile(sqlFile); err != nil {
			fmt.Printf("WARN: %s\n", err.Error())
			continue
		}
		if _, err = db.ExecContext(c.Context, string(b)); err != nil {
			fmt.Printf("WARN: %s\n", err.Error())
		}
	}
	return nil
}

func reset(c *cli.Context) error {
	if err := dropDB(c); err != nil {
		return err
	}
	if err := createDB(c); err != nil {
		return err
	}
	if err := migrateDB(c); err != nil {
		return err
	}
	if err := seedDB(c); err != nil {
		return err
	}
	return nil
}

func conn(c config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=false",
		c.user, c.password, c.host, c.port, c.dbName)
}

func connNoDB(c config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/?tls=false",
		c.user, c.password, c.host, c.port)
}
