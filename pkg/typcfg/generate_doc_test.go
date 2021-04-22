package typcfg_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
	"github.com/urfave/cli/v2"
)

func TestGenerateDoc(t *testing.T) {
	var out strings.Builder
	target := "sample.md"
	c := &typcfg.Context{
		Configs: []*typcfg.Envconfig{
			{
				Fields: []*typcfg.Field{
					{Key: "APP_NAME", Default: "some-name", Required: true},
					{Key: "APP_DEBUG", Default: "false", Required: false},
				},
			},
			{
				Fields: []*typcfg.Field{
					{Key: "DB_HOST", Default: "some-host", Required: false},
					{Key: "DB_PORT", Default: "some-port", Required: true},
				},
			},
		},

		Context: &typgo.Context{
			Descriptor: &typgo.Descriptor{
				ProjectName:    "NAME",
				ProjectVersion: "VERSION",
			},
			Context: &cli.Context{
				Command: &cli.Command{},
			},
			Logger: typgo.Logger{Stdout: &out},
		},
	}
	err := typcfg.GenerateDoc(target, c)
	require.NoError(t, err)
	defer os.Remove(target)

	b, _ := ioutil.ReadFile(target)
	expected := fmt.Sprintf(`# NAME

<!-- DO NOT EDIT. This file generated due to '@envconfig' annotation -->

## Configuration List
| Field Name | Default | Required | 
|---|---|:---:|
| APP_NAME | some-name | Yes |
| APP_DEBUG | false |  |
| DB_HOST | some-host |  |
| DB_PORT | some-port | Yes |

## DotEnv example
%s

`, "```\nAPP_NAME=some-name\nAPP_DEBUG=false\nDB_HOST=some-host\nDB_PORT=some-port\n```")
	require.Equal(t, expected, string(b))

	require.Equal(t, "> Generate 'sample.md'\n", out.String())
}
