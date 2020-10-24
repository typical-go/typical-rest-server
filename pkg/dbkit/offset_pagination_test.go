package dbkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestPagination(t *testing.T) {
	testcases := []dbkit.SelectTestCase{
		{
			SelectOption: &dbkit.OffsetPagination{},
			Builder:      sq.Select("name", "version").From("sometables"),
			Expected:     "SELECT name, version FROM sometables OFFSET 0",
		},
		{
			SelectOption: &dbkit.OffsetPagination{Offset: 10, Limit: 100},
			Builder:      sq.Select("name", "version").From("sometables"),
			Expected:     "SELECT name, version FROM sometables LIMIT 100 OFFSET 10",
		},
	}

	for _, tt := range testcases {
		tt.Execute(t)
	}
}
