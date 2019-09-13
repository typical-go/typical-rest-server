package golang

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFunction(t *testing.T) {
	pogo := Function{
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
