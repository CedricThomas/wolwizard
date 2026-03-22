package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CedricThomas/console/internal/config"
	controller "github.com/CedricThomas/console/internal/controller/base"
	redisin "github.com/CedricThomas/console/internal/input/async/redis"
	"github.com/CedricThomas/console/internal/input/async/subscriptions"
	"github.com/CedricThomas/console/internal/service/wol/wol"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Init(config.RaspberryAgent)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create Redis client for caching and async operations
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

	// Initialize external dependencies
	consumer := redisin.NewRedisConsumer(redisClient)
	wolSender := wol.New()

	// Initialize controllers
	raspberryController := controller.NewRaspberryAgentController(wolSender, &cfg.RaspberryConfig)

	// Register async subscriptions
	unsubscribes, err := subscriptions.RegisterRaspberryAgent(ctx, consumer, raspberryController)
	if err != nil {
		log.Fatalf("Failed to register async subscriptions: %v", err)
	}

	// Cleanup subscriptions
	for _, unsubscribe := range unsubscribes {
		defer func() {
			if err := unsubscribe(); err != nil {
				log.Printf("Failed to unsubscribe: %v", err)
			}
		}()
	}

	log.Println("Raspberry agent listening for async commands...")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}
