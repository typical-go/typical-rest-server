package dbkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestPagination(t *testing.T) {
	testcases := []struct {
		testName string
		*dbkit.OffsetPagination
		builder       sq.SelectBuilder
		expectedQuery string
		expectedArgs  []interface{}
	}{
		{
			OffsetPagination: &dbkit.OffsetPagination{},
			builder:          sq.Select("name", "version").From("sometables"),
			expectedQuery:    "SELECT name, version FROM sometables OFFSET 0",
		},
		{
			OffsetPagination: &dbkit.OffsetPagination{Offset: 10, Limit: 100},
			builder:          sq.Select("name", "version").From("sometables"),
			expectedQuery:    "SELECT name, version FROM sometables LIMIT 100 OFFSET 10",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			query, args, _ := tt.CompileSelect(tt.builder).ToSql()
			require.Equal(t, tt.expectedQuery, query)
			require.Equal(t, tt.expectedArgs, args)
		})
	}
}
