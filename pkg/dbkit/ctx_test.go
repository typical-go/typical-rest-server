package dbkit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestErrCtx(t *testing.T) {
	ctx := dbkit.CtxWithTxo(context.Background())
	dbkit.SetErrCtx(ctx, errors.New("some-error"))
	require.EqualError(t, dbkit.ErrCtx(ctx), "some-error")
}
