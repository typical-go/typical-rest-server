package typrest_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
)

func TestRouter(t *testing.T) {
	testcases := []struct {
		TestName string
		typrest.Router
		server      typrest.Server
		expectedErr string
	}{
		{
			Router: typrest.NewRouter(func(typrest.Server) error {
				return errors.New("some-error")
			}),
			expectedErr: "some-error",
		},
		{
			Router: typrest.Routers{
				typrest.NewRouter(func(typrest.Server) error { return errors.New("some-error-1") }),
				typrest.NewRouter(func(typrest.Server) error { return errors.New("some-error-2") }),
			},
			expectedErr: "some-error-1",
		},
		{
			Router: typrest.Routers{
				typrest.NewRouter(func(typrest.Server) error { return nil }),
				typrest.NewRouter(func(typrest.Server) error { return errors.New("some-error-2") }),
			},
			expectedErr: "some-error-2",
		},

		{
			Router: typrest.Routers{
				typrest.NewRouter(func(typrest.Server) error { return nil }),
				typrest.NewRouter(func(typrest.Server) error { return nil }),
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			err := tt.SetRoute(tt.server)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
