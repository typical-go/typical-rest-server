package typredis_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
)

func TestRedis_Task_DefaultParam(t *testing.T) {
	tool := typredis.RedisTool{}
	typgo.ProjectName = "PROJECT"
	tool.Task()

	require.Equal(t, "REDIS", tool.Name)
	require.Equal(t, typredis.EnvKeysWithPrefix("REDIS"), tool.EnvKeys)
	require.Equal(t, "PROJECT-REDIS", tool.DockerName)
}
