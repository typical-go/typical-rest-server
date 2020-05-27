package dbkit_test

import (
	"errors"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestUpdateOption(t *testing.T) {
	expectedErr := errors.New("")
	expected := sq.Update("")
	selectOpt := dbkit.NewUpdateOption(func(sq.UpdateBuilder) (sq.UpdateBuilder, error) {
		return expected, expectedErr
	})
	builder, err := selectOpt.CompileUpdate(sq.Update(""))
	require.Equal(t, expected, builder)
	require.Equal(t, expectedErr, err)
}
