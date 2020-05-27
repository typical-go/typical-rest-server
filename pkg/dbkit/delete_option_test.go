package dbkit_test

import (
	"errors"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestNewDeleteOption(t *testing.T) {
	expectedErr := errors.New("")
	expected := sq.Delete("")
	deleteOpt := dbkit.NewDeleteOption(func(sq.DeleteBuilder) (sq.DeleteBuilder, error) {
		return expected, expectedErr
	})

	builder, err := deleteOpt.CompileDelete(sq.Delete(""))
	require.Equal(t, expected, builder)
	require.Equal(t, expectedErr, err)
}
