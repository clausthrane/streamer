package source

import "context"

type (
	// Stream is an open stream to be read from
	Stream interface {
		Receive(ctx context.Context) (interface{}, error)
	}
)
