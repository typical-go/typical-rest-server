package errvalid_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/errvalid"
)

func TestValidation(t *testing.T) {
	testcases := []struct {
		testName         string
		err              error
		expectedCheck    bool
		expectedMessage  string
		expectedMessages []string
	}{
		{
			testName:         "create validation error from messages",
			err:              errvalid.New("message-1", "message-2", "message-3"),
			expectedCheck:    true,
			expectedMessage:  "message-1\nmessage-2\nmessage-3",
			expectedMessages: []string{"message-1", "message-2", "message-3"},
		},
		{
			testName:         "wrap error to validaton error ",
			err:              errvalid.Wrap(errors.New("some-error")),
			expectedCheck:    true,
			expectedMessage:  "some-error",
			expectedMessages: []string{"some-error"},
		},
		{
			testName:         "nil error",
			err:              nil,
			expectedCheck:    false,
			expectedMessage:  "",
			expectedMessages: []string{},
		},
		{
			testName:         "not validation error",
			err:              errors.New("some-error"),
			expectedCheck:    false,
			expectedMessage:  "",
			expectedMessages: []string{},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			require.Equal(t, tt.expectedCheck, errvalid.Check(tt.err))
			require.Equal(t, tt.expectedMessages, errvalid.Messages(tt.err))
		})
	}

}
