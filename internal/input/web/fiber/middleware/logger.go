package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
)

// LoggerMiddleware logs HTTP requests with method, path, status, duration, and timestamp
func LoggerMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		log.Printf(
			"[%s] %s %s %d - %v",
			time.Now().Format(time.RFC3339),
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			duration,
		)

		return err
	}
}
