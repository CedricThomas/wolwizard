package main

import (
	"log"

	"github.com/CedricThomas/console/internal/boundary/in/web/fiber/router"
	"github.com/CedricThomas/console/internal/config"
	controller "github.com/CedricThomas/console/internal/controller/base"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Cannot initialize configuration: %v", err)
	}

	app := fiber.New()
	app.Use(cors.New())
	api := app.Group("/api")
	pcaController := controller.NewPCAgentController()
	router.RegisterPCAgentRoutes(api, pcaController)
	listenAddr := ":" + cfg.Port
	log.Fatal(app.Listen(listenAddr))

	log.Println("PC-Agent server started successfully")
}
