package expr

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oklog/ulid/v2"
	"orchestrator/internal/controllers/dto"
	"orchestrator/internal/db/expressions"
	"orchestrator/internal/utils"
	"orchestrator/pkg/calc"
	"strings"
)

// calculate @Summary      Добавить выражение в очередь на выполнение
// @Description  Добавить математическое выражение из чисел, знаков [+, -, *, /, (, )] в очередь на выполнение. Возвращает ULID ID при успешном запросе
// @Tags         calculate
// @Accept       json
// @Produce      json
// @Security Bearer
// @Param        body body  dto.CalculateRequest true  "Объект, содержащий в себе выражение"
// @Success      200  {object}  dto.CalculateResponse
// @Success      201  {object}  dto.CalculateResponse
// @Failure      400  {object}  dto.ApiError
// @Failure      403  {object}  dto.ApiError
// @Failure      422  {object}  dto.ApiError
// @Failure      500  {object}  dto.ApiError
// @Router       /calculate [post]
func (ctl *Controller) calculate(ctx fiber.Ctx) error {
	var body dto.CalculateRequest
	if err := ctx.Bind().JSON(&body); err != nil {
		return utils.SendError(ctx, dto.InvalidData, fiber.StatusUnprocessableEntity)
	}

	body.Expression = strings.ReplaceAll(body.Expression, " ", "")
	body.Expression = strings.ReplaceAll(body.Expression, ",", ".")

	expr, err := ctl.exprRepo.GetByExpr(ctx.Context(), body.Expression)
	if err == nil {
		return ctx.Status(fiber.StatusOK).JSON(&dto.CalculateResponse{Id: expr.ID})
	}

	id := ulid.Make().String()
	tasks, err := calc.ParseExpression(body.Expression)
	if err != nil {
		return utils.SendError(ctx, dto.InvalidData, fiber.StatusBadRequest)
	}

	if err := ctl.exprRepo.Create(ctx.Context(), expressions.CreateParams{
		ID:       id,
		Expr:     body.Expression,
		Res:      0,
		Finished: false,
		Error:    false,
	}); err != nil {
		return utils.SendError(ctx, err.Error(), fiber.StatusInternalServerError)
	}

	for i, task := range tasks {
		if i == len(tasks)-1 {
			task.ID = id
		}
		if err := ctl.tasks.SetTask(ctx.Context(), &task); err != nil {
			return utils.SendError(ctx, err.Error(), fiber.StatusInternalServerError)
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.CalculateResponse{
		Id: id,
	})
}
