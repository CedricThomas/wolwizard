package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CedricThomas/console/internal/boundary/in/async"
	redisin "github.com/CedricThomas/console/internal/boundary/in/async/redis"
	"github.com/CedricThomas/console/internal/boundary/out/wol/wol"
	"github.com/CedricThomas/console/internal/config"
	controller "github.com/CedricThomas/console/internal/controller/base"
	asyncdomain "github.com/CedricThomas/console/internal/domain/async"
)

const bootChannel = "boot_commands"

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Parse MAC address
	if cfg.ServerMACAddress == nil {
		log.Fatal("Missing required SERVER_MAC_ADDRESS in environment")
	}
	if cfg.ServerNetworkAddress == nil {
		log.Fatal("Missing required SERVER_NETWORK_ADDRESS in environment")
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

	// Initialize external dependencies (Redis consumer and wake on lan sender)
	consumer := redisin.NewRedisConsumer(redisClient)
	wolSender := wol.New()

	controller := controller.NewRaspberryAgentController(wolSender, cfg)

	// Subscribe to boot commands and send WoL
	unsubscribe, err := async.Subscribe(ctx, consumer, asyncdomain.BootChannel, controller.ExecuteBootMessage)
	if err != nil {
		log.Fatalf("Failed to subscribe to boot channel: %v", err)
	}
	defer unsubscribe()

	log.Println("Raspberry agent listening for boot commands...")
	log.Printf("Target MAC: %s, Network: %s", cfg.ServerMACAddress, cfg.ServerNetworkAddress)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}
