package typredis_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
)

func TestModule(t *testing.T) {
	m := typredis.Module()
	require.True(t, typiobj.IsProvider(m))
	require.True(t, typiobj.IsDestroyer(m))
	require.True(t, typiobj.IsCommandLiner(m))
	require.True(t, typiobj.IsConfigurer(m))
}
