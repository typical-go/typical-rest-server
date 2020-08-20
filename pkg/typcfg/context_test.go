package typcfg_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typannot"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
)

func TestCreateField(t *testing.T) {
	testnames := []struct {
		TestName string
		Prefix   string
		Field    *typannot.Field
		Expected *typcfg.Field
	}{
		{
			Prefix:   "APP",
			Field:    &typannot.Field{Name: "Address"},
			Expected: &typcfg.Field{Key: "APP_ADDRESS"},
		},
		{
			Prefix: "APP",
			Field: &typannot.Field{
				Name:      "some-name",
				StructTag: reflect.StructTag(`envconfig:"ADDRESS" default:"some-address" required:"true"`),
			},
			Expected: &typcfg.Field{Key: "APP_ADDRESS", Default: "some-address", Required: true},
		},
	}
	for _, tt := range testnames {
		t.Run(tt.TestName, func(t *testing.T) {
			require.Equal(t, tt.Expected, typcfg.CreateField(tt.Prefix, tt.Field))
		})
	}
}
