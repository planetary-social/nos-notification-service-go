package fixtures

import (
	"context"
	"testing"
)

func Context(tb testing.TB) context.Context {
	ctx, cancelFunc := context.WithCancel(context.Background())
	tb.Cleanup(cancelFunc)
	return ctx
}
