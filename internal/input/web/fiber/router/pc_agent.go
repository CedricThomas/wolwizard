package router

import (
	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/input/web/fiber/handlers"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func RegisterPCAgentRoutes(app fiber.Router, registerCtrl controller.Register) {
	app.Use(cors.New())
	app.Post("/auth/register", handlers.Register(registerCtrl))
}
