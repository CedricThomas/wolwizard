package main

import (
	"context"
	"log"

	redisasync "github.com/CedricThomas/console/internal/boundary/out/async/redis"
	rediskeystore "github.com/CedricThomas/console/internal/boundary/out/keystore/redis"
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

	publisher := redisasync.NewRedisPublisher(redisClient)

	_ = publisher

	keystore := rediskeystore.NewRedisKeystore(redisClient)

	_ = keystore

	log.Println("Redis client initialized successfully")
}
