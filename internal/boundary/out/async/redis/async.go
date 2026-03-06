package redis

import (
	"context"

	"github.com/CedricThomas/console/internal/boundary/out/async"
	"github.com/redis/go-redis/v9"
)

// RedisPublisher implements the Publisher interface using Redis
type redisPublisher struct {
	client *redis.Client
}

// NewRedisPublisher creates a new instance of RedisPublisher with the provided Redis client
func NewRedisPublisher(client *redis.Client) (async.Publisher, error) {
	return &redisPublisher{
		client: client,
	}, nil
}

// Publish sends a message to the specified channel using Redis
func (rp *redisPublisher) Publish(ctx context.Context, channel string, message any) error {
	return rp.client.Publish(ctx, channel, message).Err()
}
