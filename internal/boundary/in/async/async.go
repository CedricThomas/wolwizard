package async

import (
	"context"
)

// Callback is a function that processes a message from a channel
type Callback func(ctx context.Context, message string) error

// Consumer defines the interface for consuming messages from channels
type Consumer interface {
	// Subscribe registers a callback for the given channel
	// Returns an error if the subscription fails
	Subscribe(ctx context.Context, channel string, callback Callback) (func() error, error)
}
