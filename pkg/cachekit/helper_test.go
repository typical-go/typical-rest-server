package cachekit_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/typical-go/typical-rest-server/pkg/cachekit"
	"github.com/stretchr/testify/require"
)

func TestNotModifiedError(t *testing.T) {
	testcases := []struct {
		desc     string
		err      error
		expected bool
	}{
		{
			desc:     "predefined no-modified error",
			err:      cachekit.ErrNotModified,
			expected: true,
		},
		{
			desc:     "predefined no-modified error with prefix",
			err:      fmt.Errorf("Prefix: %w", cachekit.ErrNotModified),
			expected: true,
		},
		{
			desc:     "random error",
			err:      errors.New("random-error"),
			expected: false,
		},
		{
			desc:     "nil error",
			err:      nil,
			expected: false,
		},
	}
	for _, tt := range testcases {
		require.Equal(t, tt.expected, cachekit.NotModifiedError(tt.err), tt.desc)
	}
}
