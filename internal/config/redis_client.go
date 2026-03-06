package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client using the provided Config
func NewRedisClient(ctx context.Context, cfg *Config) (*redis.Client, error) {
	redisOpts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %v", err)
	}

	// Create a new Redis client with the provided configuration
	rdb := redis.NewClient(redisOpts)

	// Ping the Redis server to check if it's reachable
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %v", err)
	}

	return rdb, nil
}
