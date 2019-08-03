package readme

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadmeRecipe(t *testing.T) {

	recipe := ReadmeRecipe{
		Title:       "some-title",
		Description: "some-descrption",
		Sections: []SectionPogo{
			{Title: "section-title1", Content: "section-content"},
			{Title: "section-title1", Content: "{{.Field1}}", Data: struct{ Field1 string }{"some-field1"}},
		},
	}

	var builder strings.Builder
	recipe.Write(&builder)

	require.Equal(t, `# some-title

some-descrption

## section-title1

section-content

## section-title1

some-field1

`, builder.String())
}
