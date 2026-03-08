package router

import (
	"github.com/CedricThomas/console/internal/boundary/in/web/fiber/handlers"
	"github.com/CedricThomas/console/internal/controller"
	"github.com/gofiber/fiber/v3"
)

// RegisterPCAgentRoutes is the fiber web Router for the PCAgent controller
func RegisterPCAgentRoutes(app fiber.Router, controller controller.PCAgent) {
	app.Post("/shutdown", handlers.ShutdownCurrentHost(controller))
}
