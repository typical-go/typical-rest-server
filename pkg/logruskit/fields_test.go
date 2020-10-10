package logruskit_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/logruskit"
)

func TestPutField(t *testing.T) {
	ctx := context.Background()
	logruskit.PutField(&ctx, "key1", "value1")
	require.Equal(t, logrus.Fields{"key1": "value1"}, logruskit.GetFields(ctx))
	logruskit.PutField(&ctx, "key2", 9999)
	require.Equal(t, logrus.Fields{"key1": "value1", "key2": 9999}, logruskit.GetFields(ctx))
}
