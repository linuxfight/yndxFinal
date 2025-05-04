package utils

import "github.com/gofiber/fiber/v3"

func SendError(ctx fiber.Ctx, message string, status int) error {
	return ctx.Status(status).JSON(fiber.Error{Message: message, Code: status})
}
