package recoverer

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"runtime/debug"
)

func New() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				log.Info(string(stack))
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
