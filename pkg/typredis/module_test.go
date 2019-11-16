package typredis_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typcli"
	"github.com/typical-go/typical-go/pkg/typmod"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
)

func TestModule(t *testing.T) {
	m := typredis.Module()
	require.True(t, typmod.IsProvider(m))
	require.True(t, typmod.IsDestroyer(m))
	require.True(t, typcli.IsCommander(m))
	require.True(t, typmod.IsConfigurer(m))
}
