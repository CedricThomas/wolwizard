package router

import (
	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/input/web/fiber/handlers"
	"github.com/CedricThomas/console/internal/input/web/fiber/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// RegisterWebRoutes configures and registers all web routes
func RegisterWebRoutes(app fiber.Router, webCtrl controller.Web) {
	// Configure CORS middleware
	app.Use(cors.New())

	// Configure auth middleware (authCtrl interface is embedded in webCtrl)
	authMiddleware := middleware.AuthMiddleware(webCtrl)

	// Public routes
	public := app.Group("")
	public.Post("/auth/login", handlers.Login(webCtrl))
	public.Get("/auth/verify", handlers.Verify(webCtrl))

	// API routes (require authentication)
	api := app.Group("/api")
	api.Use(authMiddleware)
	api.Post("/boot", handlers.BootSelectedOS(webCtrl))
	api.Post("/shutdown", handlers.ShutdownHandler(webCtrl))

	// Root route - serve index.html
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendFile("./static/index.html")
	})

	// Static files route
	app.Get("/static/*", func(c fiber.Ctx) error {
		filePath := "./static/" + c.Params("*", "")
		return c.SendFile(filePath)
	})
}
