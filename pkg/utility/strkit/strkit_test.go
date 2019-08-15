package strkit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToInt64(t *testing.T) {
	testcase := []struct {
		s   string
		i   int64
		err string
	}{
		{"123", int64(123), ""},
		{"abc", 0, "strconv.ParseInt: parsing \"abc\": invalid syntax"},
		{"", 0, "strconv.ParseInt: parsing \"\": invalid syntax"},
	}

	for _, tt := range testcase {
		i, err := ToInt64(tt.s)
		if tt.err != "" {
			require.EqualError(t, err, tt.err)
		}
		require.Equal(t, tt.i, i)

	}
}
