package keystore

import (
	"github.com/arzad/console/internal/boundary/out/redis_db/redis"
	"time"
)

type keystore interface {
	Get(key string) (string, error)
	Set(key string, string) error
	Delete(key string) error
	Exists(key string) (bool, error)
	Keys(pattern string) ([]string, error)
	Close() error
}

// New creates a new Redis database client
func New() *redis.DB {
	return redis.New()
}
