package sqkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
)

func TestPagination(t *testing.T) {
	testcases := []struct {
		testName string
		*sqkit.OffsetPagination
		builder       sq.SelectBuilder
		expectedQuery string
		expectedArgs  []interface{}
	}{
		{
			OffsetPagination: &sqkit.OffsetPagination{},
			builder:          sq.Select("name", "version").From("sometables"),
			expectedQuery:    "SELECT name, version FROM sometables",
		},
		{
			OffsetPagination: &sqkit.OffsetPagination{Offset: 10, Limit: 100},
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
