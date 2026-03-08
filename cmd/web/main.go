package main

import (
	"context"
	"log"

	"github.com/CedricThomas/console/internal/boundary/in/web/fiber/router"
	redisasync "github.com/CedricThomas/console/internal/boundary/out/async/redis"
	rediskeystore "github.com/CedricThomas/console/internal/boundary/out/keystore/redis"
	"github.com/CedricThomas/console/internal/config"
	controller "github.com/CedricThomas/console/internal/controller/base"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
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
	keystore := rediskeystore.NewRedisKeystore(redisClient)

	webController := controller.NewWebController(publisher, keystore)

	app := fiber.New()
	app.Use(cors.New())
	api := app.Group("/api")
	router.RegisterWebRoutes(api, webController)
	listenAddr := ":" + cfg.WebPort
	log.Fatal(app.Listen(listenAddr))

	log.Println("Redis client initialized successfully")
}
