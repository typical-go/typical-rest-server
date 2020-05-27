package dbkit

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
)

type (
	// DeleteOption to compile delete query
	DeleteOption interface {
		CompileDelete(sq.DeleteBuilder) (sq.DeleteBuilder, error)
	}
	// DeleteTestCase for testing purpose
	DeleteTestCase struct {
		TestName string
		DeleteOption
		Builder      sq.DeleteBuilder
		ExpectedErr  string
		Expected     string
		ExpectedArgs []interface{}
	}
	// CompileDeleteFn function
	CompileDeleteFn func(sq.DeleteBuilder) (sq.DeleteBuilder, error)

	deleteOptionImpl struct {
		fn CompileDeleteFn
	}
)

// NewDeleteOption return new instance of DeleteOption
func NewDeleteOption(fn CompileDeleteFn) DeleteOption {
	return &deleteOptionImpl{
		fn: fn,
	}
}

func (s *deleteOptionImpl) CompileDelete(b sq.DeleteBuilder) (sq.DeleteBuilder, error) {
	return s.fn(b)
}

//
// DeleteTestCase
//

// Execute test
func (tt *DeleteTestCase) Execute(t *testing.T) {
	t.Run(tt.TestName, func(t *testing.T) {
		builder, err := tt.CompileDelete(tt.Builder)
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
