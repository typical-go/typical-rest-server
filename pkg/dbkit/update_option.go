package dbkit

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
)

type (
	// UpdateOption to compile update query
	UpdateOption interface {
		CompileUpdate(sq.UpdateBuilder) (sq.UpdateBuilder, error)
	}
	// UpdateTestCase to test the update
	UpdateTestCase struct {
		TestName string
		UpdateOption
		Builder      sq.UpdateBuilder
		ExpectedErr  string
		Expected     string
		ExpectedArgs []interface{}
	}
	// CompileUpdateFn function
	CompileUpdateFn  func(sq.UpdateBuilder) (sq.UpdateBuilder, error)
	updateOptionImpl struct {
		fn CompileUpdateFn
	}
)

// NewUpdateOption return new instance of UpdateOption
func NewUpdateOption(fn CompileUpdateFn) UpdateOption {
	return &updateOptionImpl{fn: fn}
}

func (u *updateOptionImpl) CompileUpdate(b sq.UpdateBuilder) (sq.UpdateBuilder, error) {
	return u.fn(b)
}

//
// UpdateTestCase
//

// Execute test
func (tt *UpdateTestCase) Execute(t *testing.T) {
	t.Run(tt.TestName, func(t *testing.T) {
		builder, err := tt.CompileUpdate(tt.Builder)
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
