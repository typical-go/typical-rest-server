package typserver_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
)

func TestModule(t *testing.T) {
	m := typserver.New()
	require.True(t, typcore.IsProvider(m))
	require.True(t, typcore.IsDestroyer(m))
	require.True(t, typcore.IsConfigurer(m))
}
