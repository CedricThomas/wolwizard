package handlers

import (
	"github.com/CedricThomas/console/internal/controller"
	"github.com/gofiber/fiber/v3"
)

// ShutdownCurrentHost is handler/controller which shutdown the current host
func ShutdownCurrentHost(controller controller.PCAgent) fiber.Handler {
	return func(_ fiber.Ctx) error {
		// var requestBody entities.Book
		// err := c.Bind().Body(&requestBody)
		// if err != nil {
		// 	c.Status(http.StatusBadRequest)
		// 	return c.JSON(presenter.BookErrorResponse(err))
		// }
		// if requestBody.Author == "" || requestBody.Title == "" {
		// 	c.Status(http.StatusInternalServerError)
		// 	return c.JSON(presenter.BookErrorResponse(errors.New(
		// 		"Please specify title and author")))
		// }
		// result, err := service.InsertBook(&requestBody)
		// if err != nil {
		// 	c.Status(http.StatusInternalServerError)
		// 	return c.JSON(presenter.BookErrorResponse(err))
		// }
		// return c.JSON(presenter.BookSuccessResponse(result))
		return nil
	}
}
