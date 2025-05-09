package expr

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oklog/ulid/v2"
	"orchestrator/internal/controllers/dto"
	"orchestrator/internal/db/expressions"
	"orchestrator/internal/utils"
)

// GetById @Summary      Получить выражение по ID
// @Description  Выражение состоит из ID (ULID), Result (0 или другое число) и Status (DONE - Успешно выполнено, FAILED - Ошибка при выполнении, PROCESSING - Выполняется). Возвращает выражение при успешном запросе
// @Tags         expressions
// @Accept       json
// @Produce      json
// @Security Bearer
// @Param        id path  string true  "01JTE60CDWQ5R3QSWZBP8J6FG3"
// @Success      200  {object}  dto.GetByIdExpressionResponse
// @Failure      400  {object}  dto.ApiError
// @Failure      403  {object}  dto.ApiError
// @Failure      404  {object}  dto.ApiError
// @Router       /expressions/{id} [get]
func (ctl *Controller) GetById(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if len(id) != 26 {
		return utils.SendError(ctx, dto.InvalidData, fiber.StatusBadRequest)
	}
	if _, err := ulid.ParseStrict(id); err != nil {
		return utils.SendError(ctx, dto.InvalidData, fiber.StatusBadRequest)
	}

	expr, err := ctl.exprRepo.GetById(ctx.Context(), id)
	if err == nil {
		if task, _ := ctl.tasks.GetTask(ctx.Context(), id); task == nil && !expr.Finished {
			expr.Finished = true
			expr.Error = true
			if err := ctl.exprRepo.Update(ctx.Context(), expressions.UpdateParams{
				Res:      0,
				Finished: true,
				Error:    true,
				ID:       id,
			}); err != nil {
				return utils.SendError(ctx, err.Error(), fiber.StatusInternalServerError)
			}
		}
		return ctx.Status(fiber.StatusOK).JSON(&dto.GetByIdExpressionResponse{Expression: dto.NewExpression(expr)})
	}

	return utils.SendError(ctx, dto.NotFound, fiber.StatusNotFound)
}
