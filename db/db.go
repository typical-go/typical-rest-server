package db

import (
	// load the driver
	_ "github.com/lib/pq"
)

// // Connect to database
// func Connect(conf config.Config) (*sql.DB, error) {
// 	return sql.Open("postgres", connectionString(conf))
// }
