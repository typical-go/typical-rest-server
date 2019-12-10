package typserver_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typobj"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
)

func TestModule(t *testing.T) {
	m := &typserver.Module{}
	require.True(t, typobj.IsProvider(m))
	require.True(t, typobj.IsDestroyer(m))
	require.True(t, typobj.IsConfigurer(m))
}
