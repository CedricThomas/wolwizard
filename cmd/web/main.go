package main

// main.go - Entry point for the web server application.
// This file initializes the configuration, Redis client, controllers, and starts the Fiber web server.

import (
	"context"
	"log"

	"github.com/CedricThomas/console/internal/config"
	controller "github.com/CedricThomas/console/internal/controller/base"
	redisin "github.com/CedricThomas/console/internal/input/async/redis"
	"github.com/CedricThomas/console/internal/input/async/subscriptions"
	"github.com/CedricThomas/console/internal/input/web/fiber/router"
	redisasync "github.com/CedricThomas/console/internal/service/async/redis"
	rediskeystore "github.com/CedricThomas/console/internal/service/keystore/redis"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	// Initialize application context
	ctx := context.Background()

	// Load application configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Cannot initialize configuration: %v", err)
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

	// Initialize external dependencies (Redis publisher and keystore)
	publisher := redisasync.NewRedisPublisher(redisClient)
	keystore := rediskeystore.NewRedisKeystore(redisClient)
	consumer := redisin.NewRedisConsumer(redisClient)

	// Initialize the web controller with dependencies
	webController := controller.NewWebController(publisher, keystore)

	// Register async subscriptions
	unsubscribes, err := subscriptions.RegisterWeb(ctx, consumer, webController)
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

	// Configure and start the Fiber web server
	app := fiber.New()
	app.Use(cors.New()) // Enable CORS for cross-origin requests

	// Define API routes
	api := app.Group("/api")
	router.RegisterWebRoutes(api, webController)

	// Serve index.html for root route
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendFile("./static/index.html")
	})

	// Serve static files from /static route
	app.Get("/static/*", func(c fiber.Ctx) error {
		filePath := "./static/" + c.Params("*", "")
		return c.SendFile(filePath)
	})

	// Start the server on the configured port
	listenAddr := ":" + cfg.Port
	log.Println("Starting web server on", listenAddr)
	log.Fatal(app.Listen(listenAddr))
}
