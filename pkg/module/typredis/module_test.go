package typredis_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tiket/TIX-COMMON-GO/tixredis"
)

func TestModule(t *testing.T) {
	module := tixredis.Module()
	require.NotNil(t, module)
}
