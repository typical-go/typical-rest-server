package typpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
)

func TestModule(t *testing.T) {
	m := &typpostgres.Module{}
	require.True(t, typcore.IsProvider(m))
	require.True(t, typcore.IsDestroyer(m))
	require.True(t, typcore.IsPreparer(m))
	require.True(t, typcore.IsBuildCommander(m))
	require.True(t, typcore.IsConfigurer(m))
}
