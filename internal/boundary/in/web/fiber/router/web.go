package router

import (
	"github.com/CedricThomas/console/internal/boundary/in/web/fiber/handlers"
	"github.com/CedricThomas/console/internal/controller"
	"github.com/gofiber/fiber/v3"
)

// RegisterWebRoutes is the fiber web Router for the Web controller
func RegisterWebRoutes(app fiber.Router, controller controller.Web) {
	app.Post("/boot", handlers.BootSelectedOS(controller))
}
