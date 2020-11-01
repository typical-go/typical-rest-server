package sqkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
)

func TestNewDeleteOption(t *testing.T) {
	expected := sq.Delete("")
	deleteOpt := sqkit.NewDeleteOption(func(sq.DeleteBuilder) sq.DeleteBuilder {
		return expected
	})
	require.Equal(t, expected, deleteOpt.CompileDelete(sq.Delete("")))
}
