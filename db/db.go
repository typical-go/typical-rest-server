package db

import (
	"fmt"

	"github.com/urfave/cli"
)

// Create database
func Create(c *cli.Context) (err error) {
	fmt.Println("Create Database")
	return
}

// Drop database
func Drop(c *cli.Context) (err error) {
	fmt.Println("Drop Database")
	return
}

// Migrate database
func Migrate(c *cli.Context) (err error) {
	fmt.Println("Migrate Database")
	return
}

// Rollback database
func Rollback(c *cli.Context) (err error) {
	fmt.Println("Rollback Database")
	return
}
