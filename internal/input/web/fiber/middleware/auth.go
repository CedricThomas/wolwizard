package middleware

import (
	"errors"
	"net/http"

	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/input/web/presenters"
	"github.com/gofiber/fiber/v3"
)

// AuthMiddleware validates JWT tokens in the Authorization header
func AuthMiddleware(authCtrl controller.Auth) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(presenters.AuthError(errors.New("missing authorization header")))
		}

		token := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}

		username, err := authCtrl.ValidateToken(c.Context(), token)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(presenters.AuthError(err))
		}

		c.Locals("username", username)
		return c.Next()
	}
}
