package typdb_test

import (
	"database/sql"

	"github.com/golang-migrate/migrate"
	"github.com/typical-go/typical-rest-server/pkg/typdb"
)

type mockDBConn struct{}

var _ typdb.DBConn = (*mockDBConn)(nil)

func MockDBConn() *mockDBConn {
	return &mockDBConn{}
}

func (*mockDBConn) Connect(*typdb.Config) (*sql.DB, error) {
	return nil, nil
}
func (*mockDBConn) ConnectAdmin(*typdb.Config) (*sql.DB, error) {
	return nil, nil
}
func (*mockDBConn) Migrate(src string, cfg *typdb.Config) (*migrate.Migrate, error) {
	return nil, nil
}
