package dbtool

import (
	"fmt"
	"os"
	"time"
)

type (
	// Config config for postgres tool
	Config struct {
		DBName string
		DBUser string
		DBPass string
		Host   string
		Port   string
	}
)

// CreateMigrationFile createa migration file
func CreateMigrationFile(migrationSrc, name string) {
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
