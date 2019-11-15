package typserver_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typimodule"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
)

func TestModule(t *testing.T) {
	m := typserver.Module()
	require.True(t, typimodule.IsProvider(m))
	require.True(t, typimodule.IsDestroyer(m))
	require.True(t, typimodule.IsConfigurer(m))
}
