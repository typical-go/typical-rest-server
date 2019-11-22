package typserver_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typmodule"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
)

func TestModule(t *testing.T) {
	m := typserver.Module()
	require.True(t, typmodule.IsProvider(m))
	require.True(t, typmodule.IsDestroyer(m))
	require.True(t, typcfg.IsConfigurer(m))
}
