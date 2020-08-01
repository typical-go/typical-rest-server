package typdocker_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
)

func TestComposer(t *testing.T) {
	expectedErr := errors.New("some-error")
	expectedRecipe := &typdocker.Recipe{}

	composer := typdocker.NewCompose(func() (*typdocker.Recipe, error) {
		return expectedRecipe, expectedErr
	})

	recipe, err := composer.ComposeV3()
	require.Equal(t, expectedRecipe, recipe)
	require.Equal(t, expectedErr, err)
}
