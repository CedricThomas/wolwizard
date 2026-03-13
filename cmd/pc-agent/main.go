package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	redisin "github.com/CedricThomas/console/internal/boundary/in/async/redis"
	"github.com/CedricThomas/console/internal/boundary/in/async/subscriptions"
	"github.com/CedricThomas/console/internal/boundary/out/command"
	"github.com/CedricThomas/console/internal/boundary/out/command/linux"
	"github.com/CedricThomas/console/internal/boundary/out/command/windows"
	"github.com/CedricThomas/console/internal/config"
	controller "github.com/CedricThomas/console/internal/controller/base"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.New()
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

	// Initialize command executor based on platform
	var executor command.CommandExecutor
	switch runtime.GOOS {
	case "linux":
		executor = linux.New()
	case "windows":
		executor = windows.New()
	default:
		log.Fatalf("Unsupported operating system: %s", runtime.GOOS)
	}

	// Initialize controllers
	pcAgentController := controller.NewPCAgentController(executor)

	// Register async subscriptions
	unsubscribes, err := subscriptions.RegisterPCAgent(ctx, consumer, pcAgentController)
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
	log.Println("PC agent listening for async commands...")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}
