package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	RedisURL string `env:"REDIS_URL,required"`
}

// New creates a new Config instance with values from the environment
func New() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to process env vars: %v", err)
	}
	return &cfg, nil
}
