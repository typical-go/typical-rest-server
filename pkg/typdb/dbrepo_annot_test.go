package typdb_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typdb"
)

func TestGetDest(t *testing.T) {
	testcases := []struct {
		Name     string
		Path     string
		Expected string
	}{
		{Path: "internal/app/data_access/mysql/a.go", Expected: "internal/generated/dbrepo/app/data_access/mysql_repo"},
		{Path: "internal/app/data_access/mysql_entity/a.go", Expected: "internal/generated/dbrepo"},
		{Path: "internal/app/entity/a.go", Expected: "internal/generated/dbrepo"},
	}
	for _, tt := range testcases {
		t.Run(tt.Name, func(t *testing.T) {
			a := typdb.DBRepoAnnot{}
			require.Equal(t, tt.Expected, a.GetDest(tt.Path))
		})
	}
}
