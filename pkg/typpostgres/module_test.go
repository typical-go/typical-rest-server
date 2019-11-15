package typpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typicli"
	"github.com/typical-go/typical-go/pkg/typimodule"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
)

func TestModule(t *testing.T) {
	m := typpostgres.Module()
	require.True(t, typimodule.IsProvider(m))
	require.True(t, typimodule.IsDestroyer(m))
	require.True(t, typicli.IsCommander(m))
	require.True(t, typimodule.IsConfigurer(m))
}
