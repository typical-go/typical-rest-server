package typrest_test

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
)

func TestError(t *testing.T) {
	testcases := []struct {
		TestName string
		typrest.Error
		expectedErr     string
		expectedHTTPErr *echo.HTTPError
	}{
		{
			Error:           typrest.NewError(1, "some-message"),
			expectedErr:     "some-message",
			expectedHTTPErr: echo.NewHTTPError(1, "some-message"),
		},
		{
			Error:           typrest.NewValidErr("some-message"),
			expectedErr:     "some-message",
			expectedHTTPErr: echo.NewHTTPError(http.StatusUnprocessableEntity, "some-message"),
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			require.Equal(t, tt.expectedErr, tt.Error.Error())
			require.Equal(t, tt.expectedHTTPErr, tt.Error.HTTPError())
		})
	}
}

func TestHTTPError(t *testing.T) {
	testcases := []struct {
		TestName string
		err      error
		expected *echo.HTTPError
	}{
		{
			err:      errors.New("some-error"),
			expected: echo.NewHTTPError(http.StatusInternalServerError, "some-error"),
		},
		{
			err:      sql.ErrNoRows,
			expected: echo.NewHTTPError(http.StatusNotFound, "Not Found"),
		},
		{
			err:      typrest.NewError(99, "some-message"),
			expected: echo.NewHTTPError(99, "some-message"),
		},
		{
			err:      echo.NewHTTPError(99, "some-message"),
			expected: echo.NewHTTPError(99, "some-message"),
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			require.Equal(t, tt.expected, typrest.HTTPError(tt.err))
		})
	}
}
