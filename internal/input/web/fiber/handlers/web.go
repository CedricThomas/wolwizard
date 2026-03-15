package handlers

import (
	"net/http"

	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/domain"
	"github.com/CedricThomas/console/internal/input/web/api"
	"github.com/CedricThomas/console/internal/input/web/presenters"
	"github.com/gofiber/fiber/v3"
)

func BootSelectedOS(controller controller.Web) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req api.BootRequest

		if err := c.Bind().Body(&req); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.BootError(err, req))
		}

		if err := req.Validate(); err != nil {
			// custom error returned with message
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.BootError(err, req))
		}

		osName := domain.OSName(req.OSName)
		err := controller.SendAsyncBootCommand(c.Context(), osName)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.BootError(err, req))
		}

		return c.JSON(presenters.BootSuccess(req))
	}
}

func ShutdownHandler(controller controller.Web) fiber.Handler {
	return func(c fiber.Ctx) error {
		err := controller.SendAsyncShutdownCommand(c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ShutdownError(err))
		}

		return c.JSON(presenters.ShutdownSuccess())
	}
}
