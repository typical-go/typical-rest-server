package releaser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanMessage(t *testing.T) {
	testcases := []struct {
		message string
		cleaned string
	}{
		{message: "    abcde    \n", cleaned: "abcde"},
		{message: "some message\n\nCo-Authored-By: xx <xx@users.noreply.github.com>", cleaned: "some message"},
	}
	for _, tt := range testcases {
		require.Equal(t, tt.cleaned, clean(tt.message))
	}
}

func TestMessage(t *testing.T) {
	testcases := []struct {
		changelog string
		message   string
	}{
		{
			"5378feb rename versioning to tagging and combine goos and goarch as target",
			"rename versioning to tagging and combine goos and goarch as target",
		},
		{"5378feb ", ""},
		{"5378feb", ""},
		{"", ""},
	}
	for _, tt := range testcases {
		require.Equal(t, tt.message, message(tt.changelog))
	}
}
