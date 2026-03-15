package handlers

import (
	"errors"
	"net/http"

	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/input/web/api"
	"github.com/CedricThomas/console/internal/input/web/presenters"
	"github.com/gofiber/fiber/v3"
)

func Login(webCtrl controller.Web) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req api.LoginRequest

		if err := c.Bind().Body(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.AuthError(err))
		}

		authenticated, err := webCtrl.CheckAuth(c.Context(), req.Username, req.Password)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.AuthError(err))
		}
		if !authenticated {
			return c.Status(http.StatusUnauthorized).JSON(presenters.AuthError(errors.New("authentication failed")))
		}

		token, err := webCtrl.GenerateToken(c.Context(), req.Username)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.AuthError(err))
		}

		return c.JSON(presenters.AuthSuccess(token))
	}
}

func Verify(webCtrl controller.Web) fiber.Handler {
	return func(c fiber.Ctx) error {
		username, ok := c.Locals("username").(string)
		if !ok || username == "" {
			return c.Status(http.StatusUnauthorized).JSON(presenters.AuthError(errors.New("invalid token")))
		}

		return c.JSON(presenters.AuthVerifySuccess())
	}
}
