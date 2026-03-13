package redis

import (
	"context"
	"errors"
	"log"

	"github.com/CedricThomas/console/internal/input/async"
	"github.com/redis/go-redis/v9"
)

// redisConsumer implements the Consumer interface using go-redis
type redisConsumer struct {
	client *redis.Client
}

// NewRedisConsumer creates a new RedisConsumer
func NewRedisConsumer(client *redis.Client) async.Consumer {
	return &redisConsumer{client: client}
}

// Subscribe subscribes to a Redis channel and starts a goroutine to handle messages
func (r *redisConsumer) Subscribe(ctx context.Context, channel string, callback async.Callback) (func() error, error) {
	if r.client == nil {
		return nil, errors.New("redis client is nil")
	}

	pubsub := r.client.Subscribe(ctx, channel)

	// Wait for confirmation that subscription is created
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, err
	}

	// Channel for receiving messages from Redis
	ch := pubsub.Channel()

	// Goroutine to handle messages
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				// Call the user's callback
				if err := callback(ctx, msg.Payload); err != nil {
					log.Printf("callback error: %v", err)
				}
			case <-ctx.Done():
				if err := pubsub.Close(); err != nil {
					log.Printf("Failed to close pubsub: %v", err)
				}
				return
			}
		}
	}()

	// Return an unsubscribe function
	unsubscribe := func() error {
		if err := pubsub.Unsubscribe(ctx, channel); err != nil {
			return err
		}
		<-done // wait for goroutine to finish
		return nil
	}

	return unsubscribe, nil
}
