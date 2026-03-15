package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/CedricThomas/console/internal/config"
	controller "github.com/CedricThomas/console/internal/controller/base"
	redisin "github.com/CedricThomas/console/internal/input/async/redis"
	"github.com/CedricThomas/console/internal/input/async/subscriptions"
	"github.com/CedricThomas/console/internal/input/cron/jobs"
	cronrobfig "github.com/CedricThomas/console/internal/input/cron/robfig"
	"github.com/CedricThomas/console/internal/input/web/fiber/middleware"
	"github.com/CedricThomas/console/internal/input/web/fiber/router"
	serviceasync "github.com/CedricThomas/console/internal/service/async/redis"
	"github.com/CedricThomas/console/internal/service/command"
	"github.com/CedricThomas/console/internal/service/command/linux"
	"github.com/CedricThomas/console/internal/service/command/windows"
	rediskeystore "github.com/CedricThomas/console/internal/service/keystore/redis"
	"github.com/CedricThomas/console/internal/service/metrics"
	metricslinux "github.com/CedricThomas/console/internal/service/metrics/linux"
	metricswindows "github.com/CedricThomas/console/internal/service/metrics/windows"
	"github.com/CedricThomas/console/internal/service/token/jwt"
	"github.com/CedricThomas/console/internal/usecase/auth/base"
	"github.com/gofiber/fiber/v3"
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
		}
	}()

	// Initialize external dependencies
	consumer := redisin.NewRedisConsumer(redisClient)
	publisher := serviceasync.NewRedisPublisher(redisClient)
	cronService := cronrobfig.NewRobfigScheduler()
	defer func() {
		if err := cronService.Stop(); err != nil {
			log.Printf("Failed to stop cron: %v", err)
		}
	}()

	// Initialize command executor and metrics collector based on the OS
	var executor command.CommandExecutor
	var collector metrics.Collector
	switch runtime.GOOS {
	case "linux":
		collector = metricslinux.New()
		executor = linux.New()
	case "windows":
		executor = windows.New()
		collector = metricswindows.New()
	default:
		log.Fatalf("Unsupported operating system: %s", runtime.GOOS)
	}

	// Initialize controllers
	authCtrl := base.New(rediskeystore.NewRedisKeystore(redisClient), jwt.New(cfg.JWTSecret, cfg.JWTExpirySeconds))
	pcAgentController := controller.NewPCAgentController(executor, collector, publisher, authCtrl)

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

	// Register periodic jobs
	if err := jobs.RegisterPCAgent(ctx, cronService, pcAgentController, cfg); err != nil {
		log.Fatalf("Failed to register periodic jobs: %v", err)
	}

	cronService.Start()

	log.Println("PC agent listening for async commands...")

	// Start http server for registration
	httpServer := fiber.New()

	router.RegisterPCAgentRoutes(httpServer, pcAgentController)

	// Register logging middleware
	httpServer.Use(middleware.LoggerMiddleware())

	listenAddr := "0.0.0.0:" + cfg.Port
	log.Printf("Starting pc-agent web server on %s", listenAddr)

	go func() {
		if err := httpServer.Listen(listenAddr); err != nil {
			log.Printf("Failed to start web server: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}
