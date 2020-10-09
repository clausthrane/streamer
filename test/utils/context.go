package utils

import "context"

func newContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

// WithContext executes the given function providing it with a
// context that is cancelled after execution.
//
// this helps ensure that no resources are leaked post execution
func WithContext(f func(context.Context)) {
	ctx, cancel := newContext()
	defer cancel()
	f(ctx)
}
