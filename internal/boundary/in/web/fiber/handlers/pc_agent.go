package handlers

import (
	"net/http"

	"github.com/CedricThomas/console/internal/boundary/in/web/presenters"
	"github.com/CedricThomas/console/internal/controller"
	"github.com/gofiber/fiber/v3"
)

// ShutdownCurrentHost is handler/controller which shutdown the current host
func ShutdownCurrentHost(controller controller.PCAgent) fiber.Handler {
	return func(c fiber.Ctx) error {
		err := controller.ShutdowncurrentHost(c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ShutdownError(err))
		}

		return c.JSON(presenters.ShutdownSuccess())
	}
}
