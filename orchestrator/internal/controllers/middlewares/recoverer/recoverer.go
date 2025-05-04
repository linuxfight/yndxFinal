package recoverer

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
)

func New() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				_ = c.Status(fiber.StatusInternalServerError).JSON(
					&fiber.Error{
						Message: fmt.Sprintf("%v", err),
						Code:    fiber.StatusInternalServerError},
				)
			}
		}()

		return c.Next()
	}
}
