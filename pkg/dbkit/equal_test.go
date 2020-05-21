package dbkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestEqualOption_Select(t *testing.T) {
	testcases := []SelectTestCase{
		{
			SelectOption: dbkit.Equal("", ""),
			Builder:      sq.Select("name", "version").From("some-table"),
			ExpectedErr:  "equal: column is missing",
		},
		{
			SelectOption: dbkit.Equal("name", "dummy-name"),
			Builder:      sq.Select("name", "version").From("some-table"),
			Expected:     "SELECT name, version FROM some-table WHERE name = ?",
			ExpectedArgs: []interface{}{"dummy-name"},
		},
	}

	for _, tt := range testcases {
		tt.Execute(t)
	}
}

func TestEqualOption_Update(t *testing.T) {
	testcases := []UpdateTestCase{
		{
			UpdateOption: dbkit.Equal("", ""),
			Builder:      sq.Update("some-table"),
			ExpectedErr:  "equal: column is missing",
		},
		{
			UpdateOption: dbkit.Equal("name", "dummy-name"),
			Builder:      sq.Update("some-table").Set("column", "column-value"),
			Expected:     "UPDATE some-table SET column = ? WHERE name = ?",
			ExpectedArgs: []interface{}{"column-value", "dummy-name"},
		},
	}

	for _, tt := range testcases {
		tt.Execute(t)
	}

}

func TestEqualOption_Delete(t *testing.T) {
	testcases := []DeleteTestCase{
		{
			DeleteOption: dbkit.Equal("", ""),
			Builder:      sq.Delete("some-table"),
			ExpectedErr:  "equal: column is missing",
		},
		{
			DeleteOption: dbkit.Equal("name", "dummy-name"),
			Builder:      sq.Delete("some-table"),
			Expected:     "DELETE FROM some-table WHERE name = ?",
			ExpectedArgs: []interface{}{"dummy-name"},
		},
	}

	for _, tt := range testcases {
		tt.Execute(t)
	}
}
