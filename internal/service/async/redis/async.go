package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CedricThomas/console/internal/service/async"
	"github.com/redis/go-redis/v9"
)

// RedisPublisher implements the Publisher interface using Redis
type redisPublisher struct {
	client *redis.Client
}

// NewRedisPublisher creates a new instance of RedisPublisher with the provided Redis client
func NewRedisPublisher(client *redis.Client) async.Publisher {
	return &redisPublisher{
		client: client,
	}
}

// Publish sends a message to the specified channel using Redis
func (rp *redisPublisher) Publish(ctx context.Context, channel string, message any) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal boot message: %w", err)
	}
	return rp.client.Publish(ctx, channel, data).Err()
}
