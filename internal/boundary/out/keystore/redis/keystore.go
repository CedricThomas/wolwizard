package redis

import (
	"context"
	"errors"

	"github.com/CedricThomas/console/internal/boundary/out/keystore"

	"github.com/redis/go-redis/v9"
)

type redisKeystore struct {
	client *redis.Client
}

func New(client *redis.Client) keystore.Keystore {
	return &redisKeystore{
		client: client,
	}
}

func (r *redisKeystore) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func (r *redisKeystore) Set(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *redisKeystore) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}

func (r *redisKeystore) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *redisKeystore) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, _, err := r.client.Scan(ctx, 0, pattern, 0).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}
