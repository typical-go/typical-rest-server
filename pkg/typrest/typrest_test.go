package typrest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typcore"
)

func TestModule(t *testing.T) {
	m := &Module{}
	require.True(t, typcore.IsBuildCommander(m))
}
