package asynchandler

import (
	"context"
	"sync"
	"testing"

	"github.com/clausthrane/streamer/test/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestGivenAHandlerFuncTheAdapterExecutesItAsync(t *testing.T) {
	utils.WithContext(func(ctx context.Context) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		var v interface{}
		a := NewAdapter(ctx, 1, func(ctx context.Context, in interface{}, out interface{}, connectionError error) {
			v = in
			wg.Done()
		})
		a.Wrapped()(ctx, "input", nil, nil)
		wg.Wait()
		assert.Equal(t, "input", v)
		a.Stop()
		goleak.VerifyNone(t)
	})
}
