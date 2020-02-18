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
		{opt: dbkit.Pagination(10, 100), expected: "pagination:10:100"},
		{opt: dbkit.Sort("some-col", dbkit.Desc), expected: "sort:some-col:1"},
		{opt: dbkit.Like("some-col", "some-cond"), expected: "like:some-col:some-cond"},
		{opt: dbkit.Equal("some-col", "some-cond"), expected: "equal:some-col:some-cond"},
	}

	for _, tt := range testcases {
		require.Equal(t, tt.expected, tt.opt.String())
	}
}
