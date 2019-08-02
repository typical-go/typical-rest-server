package generated

import (
	"strings"
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

	var builder strings.Builder
	pogo.Write(&builder)

	require.Equal(t, `func hello(message string,) (string,error){ 
return "world", nil
}`, builder.String())
}
