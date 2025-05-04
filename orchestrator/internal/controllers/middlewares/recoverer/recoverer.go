package recoverer

import (
	"github.com/gofiber/fiber/v3"
)

func New() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				_ = c.Status(fiber.StatusInternalServerError).JSON(
					&fiber.Error{
						Message: err.(error).Error(),
						Code:    fiber.StatusInternalServerError},
				)
			}
		}()

		return c.Next()
	}
}
