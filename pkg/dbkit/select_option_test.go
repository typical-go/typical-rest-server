package dbkit_test

import (
	"errors"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestSelectOption(t *testing.T) {
	expectedErr := errors.New("")
	expected := sq.Select("")
	selectOpt := dbkit.NewSelectOption(func(sq.SelectBuilder) (sq.SelectBuilder, error) {
		return expected, expectedErr
	})
	builder, err := selectOpt.CompileSelect(sq.Select(""))
	require.Equal(t, expected, builder)
	require.Equal(t, expectedErr, err)
}
