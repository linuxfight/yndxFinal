package solver

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oklog/ulid/v2"
	"orchestrator/internal/controllers/dto"
	"orchestrator/internal/controllers/utils"
)

// GetById @Summary      Получить выражение по ULID
// @Tags         expressions
// @Accept       json
// @Produce      json
// @Security Bearer
// @Param        id path  string true  "ULID выражения"
// @Success      200  {object}  dto.GetByIdExpressionResponse
// @Failure      404  {object}  dto.ApiError
// @Failure      422  {object}  dto.ApiError
// @Router       /expressions/{id} [get]
func (ctl *Controller) GetById(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if _, err := ulid.Parse(id); err != nil {
		return utils.SendError(ctx, dto.InvalidData, fiber.StatusUnprocessableEntity)
	}

	expr, err := ctl.exprRepo.GetById(ctx.Context(), id)
	if err == nil {
		return ctx.Status(fiber.StatusOK).JSON(&dto.GetByIdExpressionResponse{Expression: dto.NewExpression(expr)})
	}

	return utils.SendError(ctx, dto.NotFound, fiber.StatusNotFound)
}
