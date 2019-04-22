package repo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBook_Validate(t *testing.T) {
	book := Book{}
	err := book.Validate()
	require.EqualError(t, err, `Key: 'Book.Title' Error:Field validation for 'Title' failed on the 'required' tag
Key: 'Book.Author' Error:Field validation for 'Author' failed on the 'required' tag`)
}
