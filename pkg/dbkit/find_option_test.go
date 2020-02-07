package dbkit_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sq "github.com/Masterminds/squirrel"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestCreateFindOption(t *testing.T) {
	t.Run("Must return implement of FindOption interface", func(t *testing.T) {
		var _ dbkit.FindOption = dbkit.CreateFindOption()
	})
}

func TestFindOption(t *testing.T) {
	testcases := []struct {
		dbkit.FindOption
		base          sq.SelectBuilder
		expectedError string
		expected      string
	}{
		{
			FindOption: dbkit.CreateFindOption().WithPagination(dbkit.NewPagination(0, 100)),
			base:       sq.Select("some-column").From("some-table"),
			expected:   "SELECT some-column FROM some-table LIMIT 100 OFFSET 0",
		},
	}

	for _, tt := range testcases {
		base, err := tt.CompileQuery(tt.base)
		if err != nil {
			require.EqualError(t, err, tt.expectedError)
		} else {
			sql, _, _ := base.ToSql()
			require.Equal(t, tt.expected, sql)
		}

	}

}
