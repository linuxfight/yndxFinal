package expr

import (
	"github.com/gofiber/fiber/v3"
	"orchestrator/internal/controllers/dto"
	"orchestrator/internal/db/expressions"
	"orchestrator/internal/utils"
)

// list @Summary      Получить весь список выражений
// @Description  Выражение состоит из ID (ULID), Result (0 или другое число) и Status (DONE - Успешно выполнено, FAILED - Ошибка при выполнении, PROCESSING - Выполняется). Возвращает список выражений при успешном запросе
// @Tags         expressions
// @Accept       json
// @Produce      json
// @Security Bearer
// @Success      200  {object}  dto.ListAllExpressionsResponse
// @Failure      500  {object}  dto.ApiError
// @Router       /expressions [get]
func (ctl *Controller) list(ctx fiber.Ctx) error {
	exprs, err := ctl.exprRepo.GetAll(ctx.Context())
	if err != nil {
		return utils.SendError(ctx, err.Error(), fiber.StatusInternalServerError)
	}

	result := []dto.Expression{}

	for _, expr := range exprs {
		if task, _ := ctl.tasks.GetTask(ctx.Context(), expr.ID); task == nil && !expr.Finished {
			expr.Finished = true
			expr.Error = true
			if err := ctl.exprRepo.Update(ctx.Context(), expressions.UpdateParams{
				Res:      0,
				Finished: true,
				Error:    true,
				ID:       expr.ID,
			}); err != nil {
				return utils.SendError(ctx, err.Error(), fiber.StatusInternalServerError)
			}
		}

		result = append(result, dto.NewExpression(expr))
	}

	return ctx.Status(fiber.StatusOK).JSON(&dto.ListAllExpressionsResponse{Expressions: result})
}
