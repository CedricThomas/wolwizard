package async

import (
	"context"
	"encoding/json"
	"fmt"
)

// Callback is a function that processes a message from a channel
type Callback func(ctx context.Context, message string) error

// Consumer defines the interface for consuming messages from channels
type Consumer interface {
	// Subscribe registers a callback for the given channel
	// Returns an error if the subscription fails
	Subscribe(ctx context.Context, channel string, callback Callback) (func() error, error)
}

// TypedCallback is a function that processes a message from a channel
type TypedCallback[T any] func(ctx context.Context, message T) error

// Generic typed subscription helper
func Subscribe[T any](ctx context.Context, consumer Consumer, channel string, handler TypedCallback[T]) (func() error, error) {
	wrapped := func(ctx context.Context, msg string) error {
		var data T
		if err := json.Unmarshal([]byte(msg), &data); err != nil {
			return fmt.Errorf("invalid unmarshaling on consumption: %v", err)
		}
		return handler(ctx, data)
	}
	return consumer.Subscribe(ctx, channel, wrapped)
}
