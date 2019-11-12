package typienv_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

func TestEnvVar(t *testing.T) {
	os.Setenv("TEST_2", "value-2")
	defer os.Clearenv()
	testcases := []struct {
		envVar typienv.EnvVar
		value  string
	}{
		{typienv.EnvVar{"TEST_1", "default-1"}, "default-1"},
		{typienv.EnvVar{"TEST_2", "default-2"}, "value-2"},
	}
	for i, tt := range testcases {
		require.Equal(t, tt.value, tt.envVar.Value(), i)
	}
}
