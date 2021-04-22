package typcfg_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
)

func TestCreateAndLoadDotEnv_EnvFileExist(t *testing.T) {
	target := "some-env"
	ioutil.WriteFile(target, []byte("key1=val111\nkey2=val222"), 0777)

	defer os.Remove(target)

	var out strings.Builder
	c := &typgo.Context{Logger: typgo.Logger{Stdout: &out}}
	cc := &typcfg.Context{
		Context: c,
		Configs: []*typcfg.Envconfig{
			{
				Fields: []*typcfg.Field{
					{Key: "key1", Default: "val1"},
					{Key: "key2", Default: "val2"},
					{Key: "key3", Default: "val3"},
				},
			},
		},
	}

	typcfg.GenerateAndLoadDotEnv(target, cc)

	b, _ := ioutil.ReadFile(target)
	require.Equal(t, "key1=val111\nkey2=val222\nkey3=val3\n", string(b))
	require.Equal(t, "> New keys added in 'some-env': key3\n", out.String())
}
