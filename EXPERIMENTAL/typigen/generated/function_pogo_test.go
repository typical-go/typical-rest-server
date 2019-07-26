package generated

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFunctionPogo(t *testing.T) {
	pogo := FunctionPogo{
		Name: "hello",
		FuncParams: map[string]string{
			"message": "string",
		},
		ReturnValues: []string{"string", "error"},
		FuncBody:     `return "world", nil`,
	}

	require.Equal(t, `func hello(message string,) (string,error){ 
return "world", nil
}
`, pogo.String())
}
