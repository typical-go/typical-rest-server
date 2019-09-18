package typredis_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/module/typredis"
)

func TestModule(t *testing.T) {
	module := typredis.Module()
	require.NotNil(t, module)
}
