package dbkit_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestFindOption_String(t *testing.T) {
	testcases := []struct {
		opt      dbkit.FindOption
		expected string
	}{
		{opt: dbkit.Pagination(10, 100), expected: "pagination from 10 limit 100"},
		{opt: dbkit.Sort("some-col", dbkit.Desc), expected: "sort some-col by 1"},
		{opt: dbkit.Like("some-col", "some-cond"), expected: "where some-col like some-cond"},
		{opt: dbkit.Equal("some-col", "some-cond"), expected: "where some-col=some-cond"},
	}

	for _, tt := range testcases {
		require.Equal(t, tt.expected, tt.opt.String())
	}
}
