package handlers

import (
	"net/http"

	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/input/web/api"
	"github.com/CedricThomas/console/internal/input/web/presenters"
	"github.com/gofiber/fiber/v3"
)

func Register(registerCtrl controller.Register) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req api.RegisterRequest

		if err := c.Bind().Body(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.RegisterError(err))
		}

		if err := req.Validate(); err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.RegisterError(err))
		}

		if err := registerCtrl.CreateAccount(c.Context(), req.Username, req.Password); err != nil {
			return c.Status(http.StatusConflict).JSON(presenters.RegisterError(err))
		}

		return c.JSON(presenters.RegisterSuccess())
	}
}
