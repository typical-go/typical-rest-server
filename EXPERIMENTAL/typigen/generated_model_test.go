package typigen

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratedModel(t *testing.T) {
	model := GeneratedModel{
		PackageName: "typical",
		Imports: map[string]string{
			"fmt":                              "",
			"github.com/typical-go/typical-go": "tgo",
		},
		Structs: []StructPogo{
			{
				Name: "Config",
				Fields: []reflect.StructField{
					{Name: "Name", Type: reflect.TypeOf("something")},
				},
			},
		},

		AddConstructors: []FunctionPogo{
			{
				FuncParams:   map[string]string{"test": "string"},
				ReturnValues: []string{"string"},
				FuncBody:     `return "hello"`,
			},
		},
	}

	// fmt.Println(model.String())

	require.Equal(t, `package typical
import  "fmt"
import tgo "github.com/typical-go/typical-go"
type Config struct{
Name string
}

func init() {
Context.AddConstructor(func (test string,) (string){ 
return "hello"
}
)
}
`, model.String())
}
