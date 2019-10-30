package typserver_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
)

func TestModule(t *testing.T) {
	m := typserver.Module()
	require.True(t, typiobj.IsProvider(m))
	require.True(t, typiobj.IsDestructor(m))
	require.True(t, typiobj.IsConfigurer(m))
}
