package errkit_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/utility/errkit"
)

func TestErrors(t *testing.T) {
	var errors errkit.Errors
	errors.Add(fmt.Errorf("error1"))
	errors.Add(nil)
	errors.Add(fmt.Errorf("error2"))
	errors.Add(fmt.Errorf("error3"))

	require.Equal(t, "error1; error2; error3", errors.Error())

}
