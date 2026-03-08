package main

import (
	"context"
	"log"

	asyncredis "github.com/CedricThomas/console/internal/boundary/out/async/redis"
	"github.com/CedricThomas/console/internal/config"
)

func main() {
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Cannot initialize configuration: %v", err)
	}

	redisClient, err := config.NewRedisClient(ctx, cfg)
	if err != nil {
		log.Fatalf("Cannot initialize Redis client: %v", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Failed to close Redis client: %v", err)
			return
		}
	}()

	publisher, err := asyncredis.NewRedisPublisher(redisClient)
	if err != nil {
		log.Fatalf("Cannot initialize Redis publisher: %v", err)
	}

	_ = publisher

	log.Println("Redis client initialized successfully")
}
