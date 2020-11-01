package sqkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
)

func TestUpdateOption(t *testing.T) {
	expected := sq.Update("")
	selectOpt := sqkit.NewUpdateOption(func(sq.UpdateBuilder) sq.UpdateBuilder {
		return expected
	})

	require.Equal(t, expected, selectOpt.CompileUpdate(sq.Update("")))

}
