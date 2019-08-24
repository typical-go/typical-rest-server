package typirelease

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanMessage(t *testing.T) {
	testcase := []struct {
		message string
		cleaned string
	}{
		{message: "    abcde    \n", cleaned: "abcde"},
		{message: "some message\n\nCo-Authored-By: xx <xx@users.noreply.github.com>", cleaned: "some message"},
	}

	for _, tt := range testcase {
		result := cleanMessage(tt.message)
		require.Equal(t, tt.cleaned, result)
	}
}
