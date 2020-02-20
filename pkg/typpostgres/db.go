package typpostgres

import "database/sql"

// DB is postgres database handle
type DB struct {
	*sql.DB
}

// NewDB return new instance of DB
func NewDB(db *sql.DB) *DB {
	return &DB{
		DB: db,
	}
}
