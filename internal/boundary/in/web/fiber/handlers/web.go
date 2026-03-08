package handlers

import (
	"net/http"

	"github.com/CedricThomas/console/internal/boundary/in/web/api"
	"github.com/CedricThomas/console/internal/boundary/in/web/presenters"
	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/domain"
	"github.com/gofiber/fiber/v3"
)

func BootSelectedOS(controller controller.Web) fiber.Handler {
	return func(c fiber.Ctx) error {

		var req api.BootRequest

		if err := c.Bind().Body(&req); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.BootError(err))
		}

		if err := req.Validate(); err != nil {
			// custom error returned with message
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.BootError(err))
		}

		osName := domain.OSName(req.Name)
		err := controller.BootSelectedOS(c.Context(), osName)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.BootError(err))
		}

		return c.JSON(presenters.BootSuccess(osName))
	}
}
