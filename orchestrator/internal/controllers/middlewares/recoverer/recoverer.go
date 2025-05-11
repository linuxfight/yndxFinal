package recoverer

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"runtime/debug"
)

func New() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				_ = c.Status(fiber.StatusInternalServerError).JSON(
					&fiber.Error{
						Message: fmt.Sprintf("%v \n stack: %s", err, string(stack)),
						Code:    fiber.StatusInternalServerError},
				)
			}
		}()

		return c.Next()
	}
}
