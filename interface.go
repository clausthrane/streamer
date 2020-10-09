package streamer

import (
	"context"

	"github.com/clausthrane/streamer/source"
)

type (
	//StreamHandlerFunc is a function to be called on each event from an open stream
	StreamHandlerFunc = func(ctx context.Context, in interface{}, out interface{}, connectionError error)

	//StreamHandler holds a Handle is a function to be called on each event from an open stream
	StreamHandler interface {
		Handle(ctx context.Context, in interface{}, out interface{}, connectionError error)
	}

	// Connector creates a connection to a stream
	Connector interface {
		ID() string
		Connect(ctx context.Context) (source.Stream, error)
		StreamOpen() bool
	}
)
