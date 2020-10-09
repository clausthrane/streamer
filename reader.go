package streamer

import (
	"context"

	"github.com/clausthrane/streamer/source"
)

type (
	// A Reader is an entity which encapsulates the processing of a connection from a Connector and
	// invokes a StreamHandlerFunc with the result.
	Reader struct {
		connector  Connector
		connection source.Stream
		ctx        context.Context
		ctxCloser  context.CancelFunc
		handler    StreamHandlerFunc
	}
)

// NewReader creates a new Reader instance associated with the given context.
//
// All processing will happen in a child context and the caller can start the process
// by calling Start() and close it by cancelling the provided context or calling Stop()
// once stopped the Reader instance cannot be reused.
func NewReader(application context.Context, c Connector, h StreamHandlerFunc) *Reader {
	ctx, cancel := context.WithCancel(application)
	return &Reader{
		ctx:       ctx,
		ctxCloser: cancel,
		connector: c,
		handler:   h,
	}
}

// Active asserts whether the underlying connection is open.
// as long and the Reader has been started the stream will be read.
func (r *Reader) Active() bool {
	return r.connector.StreamOpen()
}

// Start connects the stream and initiates the async processing
func (r *Reader) Start() error {
	err := r.connect()
	if err != nil {
		return err
	}
	r.asyncReceive()
	return nil
}

func (r *Reader) connect() error {
	stream, err := r.connector.Connect(r.ctx)
	r.connection = stream
	return err
}

func (r *Reader) asyncReceive() {
	go r.syncReceive()
}

func (r *Reader) syncReceive() {
	for {
		if r.ctx.Err() != nil {
			return
		}
		r.receiveOnce()
	}
}

func (r *Reader) receiveOnce() {
	value, err := r.connection.Receive(r.ctx)
	r.handler(r.ctx, value, nil, err)
}

// Stop closes the context for this Reader and the underlying connection
func (r *Reader) Stop() {
	r.ctxCloser()
}
