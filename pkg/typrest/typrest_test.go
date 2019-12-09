package typrest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typcli"
	"github.com/typical-go/typical-go/pkg/typmodule"
)

func TestModule(t *testing.T) {
	m := &Module{}
	require.True(t, typmodule.IsProvider(m))
	require.True(t, typmodule.IsDestroyer(m))
	require.True(t, typmodule.IsProvider(m))
	require.True(t, typcli.IsBuildCommander(m))
}
