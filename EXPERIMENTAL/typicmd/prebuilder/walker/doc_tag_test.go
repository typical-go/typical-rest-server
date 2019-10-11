package walker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocTags_ParseDocTag(t *testing.T) {
	tags := ParseDocTag("[tag1] some text [tag2] then another text [tag3]")
	require.Equal(t, DocTags{"tag1", "tag2", "tag3"}, tags)
}

func TestDocTags_Contain(t *testing.T) {
	testcases := []struct {
		doctags DocTags
		tag     string
		contain bool
	}{
		{[]string{"tag1", "tag2"}, "tag1", true},
		{[]string{"tag1", "tag2"}, "tag3", false},
	}
	for _, tt := range testcases {
		result := tt.doctags.Contain(tt.tag)
		require.Equal(t, result, tt.contain)
	}
}
