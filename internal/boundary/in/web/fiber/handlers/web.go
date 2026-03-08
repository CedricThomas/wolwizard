package handlers

import (
	"github.com/CedricThomas/console/internal/controller"
	"github.com/gofiber/fiber/v3"
)

// BootWindows is handler/controller which send a message to boot the designed OS
func BootSelectedOS(controller controller.Web) fiber.Handler {
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
