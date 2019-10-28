package typpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/pkg/module/typpostgres"
)

func TestModule(t *testing.T) {
	m := typpostgres.Module()
	_, ok := m.(typictx.Constructor)
	require.True(t, ok)
	_, ok = m.(typictx.Destructor)
	require.True(t, ok)
	_, ok = m.(typictx.CommandLiner)
	require.True(t, ok)
}
