package asynchandler

import (
	"context"

	"github.com/clausthrane/streamer"
)

type (
	// Adapter wraps a given StreamHandlerFunc and buffers the original input executing the
	// buffered input on a go routine
	Adapter struct {
		cancelFunc context.CancelFunc
		ctx        context.Context
		buf        chan bufItem
		inner      streamer.StreamHandlerFunc
	}

	bufItem struct {
		ctx             context.Context
		in              interface{}
		connectionError error
	}
)

// NewAdapter creates a wrapper for the given StreamHandlerFunc
func NewAdapter(application context.Context, bufSize int, inner streamer.StreamHandlerFunc) *Adapter {
	h := newAdapter(application, bufSize, inner)
	h.asyncDequeue()
	return h
}

func newAdapter(application context.Context, bufSize int, inner streamer.StreamHandlerFunc) *Adapter {
	ctx, cancel := context.WithCancel(application)
	return &Adapter{
		ctx:        ctx,
		cancelFunc: cancel,
		buf:        make(chan bufItem, bufSize),
		inner:      inner,
	}
}

// Wrapped returns the wrapped StreamHandlerFunc
func (a *Adapter) Wrapped() streamer.StreamHandlerFunc {
	return func(ctx context.Context, in interface{}, _ interface{}, connectionError error) {
		select {
		case <-ctx.Done():
			return
		case a.buf <- bufItem{ctx, in, connectionError}:
		}
	}
}

// Stop terminates the underlying goroutine
func (a *Adapter) Stop() {
	a.cancelFunc()
}

func (a *Adapter) asyncDequeue() {
	go a.syncDequeue()
}

func (a *Adapter) syncDequeue() {
	for {
		if a.ctx.Err() != nil {
			return
		}
		a.dequeueOnce()
	}
}

func (a *Adapter) dequeueOnce() {
	select {
	case bufItem := <-a.buf:
		a.inner(bufItem.ctx, bufItem.in, nil, bufItem.connectionError)
	case <-a.ctx.Done():
		return
	}
}
