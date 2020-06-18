package dbkit

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
)

type (
	// SelectOption to compile select query
	SelectOption interface {
		CompileSelect(sq.SelectBuilder) (sq.SelectBuilder, error)
	}
	// SelectTestCase for testing puspose
	SelectTestCase struct {
		TestName string
		SelectOption
		Builder      sq.SelectBuilder
		ExpectedErr  string
		Expected     string
		ExpectedArgs []interface{}
	}
	// CompileSelectFn function
	CompileSelectFn  func(sq.SelectBuilder) (sq.SelectBuilder, error)
	selectOptionImpl struct {
		fn CompileSelectFn
	}
)

// NewSelectOption return new instance of SelectOption
func NewSelectOption(fn CompileSelectFn) SelectOption {
	return &selectOptionImpl{fn: fn}
}

func (s *selectOptionImpl) CompileSelect(b sq.SelectBuilder) (sq.SelectBuilder, error) {
	return s.fn(b)
}

//
// SelectTestCase
//

// Execute test
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
