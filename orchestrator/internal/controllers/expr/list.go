package expr

import (
	"github.com/gofiber/fiber/v3"
	"orchestrator/internal/controllers/dto"
	"orchestrator/internal/utils"
)

// list @Summary      Получить весь список выражений
// @Tags         expressions
// @Accept       json
// @Produce      json
// @Security Bearer
// @Success      200  {object}  dto.ListAllExpressionsResponse
// @Failure      500  {object}  dto.ApiError
// @Router       /expressions [get]
func (ctl *Controller) list(ctx fiber.Ctx) error {
	expressions, err := ctl.exprRepo.GetAll(ctx.Context())
	if err != nil {
		return utils.SendError(ctx, err.Error(), fiber.StatusInternalServerError)
	}

	result := []dto.Expression{}

	for _, expr := range expressions {
		result = append(result, dto.NewExpression(expr))
	}

	return ctx.Status(fiber.StatusOK).JSON(&dto.ListAllExpressionsResponse{Expressions: result})
}
