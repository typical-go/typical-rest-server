package typtool_test

import (
	"database/sql"

	"github.com/golang-migrate/migrate"
	"github.com/typical-go/typical-rest-server/pkg/typtool"
)

type mockDBConn struct{}

var _ typtool.DBConn = (*mockDBConn)(nil)

func MockDBConn() *mockDBConn {
	return &mockDBConn{}
}

func (*mockDBConn) Connect(*typtool.DBConfig) (*sql.DB, error) {
	return nil, nil
}
func (*mockDBConn) ConnectAdmin(*typtool.DBConfig) (*sql.DB, error) {
	return nil, nil
}
func (*mockDBConn) Migrate(src string, cfg *typtool.DBConfig) (*migrate.Migrate, error) {
	return nil, nil
}
