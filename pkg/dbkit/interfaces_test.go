package dbkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

type (
	SelectTestCase struct {
		TestName string
		dbkit.SelectOption
		Builder      sq.SelectBuilder
		ExpectedErr  string
		Expected     string
		ExpectedArgs []interface{}
	}
)

func (tt *SelectTestCase) Execute(t *testing.T) {
	t.Run(tt.TestName, func(t *testing.T) {
		builder, err := tt.CompileSelect(tt.Builder)
		if tt.ExpectedErr != "" {
			require.EqualError(t, err, tt.ExpectedErr)
			return
		}

		require.NoError(t, err)
		query, args, _ := builder.ToSql()
		require.Equal(t, tt.Expected, query)
		require.Equal(t, tt.ExpectedArgs, args)
	})
}
