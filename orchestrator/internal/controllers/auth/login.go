package auth

import (
	"github.com/gofiber/fiber/v3"
	"orchestrator/internal/controllers/dto"
	"orchestrator/internal/utils"
	"time"
)

// login godoc
// @Summary      Войти в аккаунт
// @Description  Войти в аккаунт с помощью логина и пароля. Возвращает JWT API Token при успешном запросе
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body  dto.AuthRequest true  "Данные для входа"
// @Success      200  {object}  dto.AuthResponse
// @Failure      400  {object}  dto.ApiError
// @Failure      401  {object}  dto.ApiError
// @Failure      404  {object}  dto.ApiError
// @Failure      422  {object}  dto.ApiError
// @Failure      500  {object}  dto.ApiError
// @Router       /login [post]
func (ctl *Controller) login(ctx fiber.Ctx) error {
	var body dto.AuthRequest
	if err := ctx.Bind().JSON(&body); err != nil {
		return utils.SendError(ctx, dto.InvalidData, fiber.StatusUnprocessableEntity)
	}

	if len(body.Login) < 1 || len(body.Password) < 1 {
		return utils.SendError(ctx, dto.InvalidData, fiber.StatusBadRequest)
	}

	user, err := ctl.usersRepo.GetByName(ctx.Context(), body.Login)
	if err != nil {
		return utils.SendError(ctx, dto.NotFound, fiber.StatusNotFound)
	}

	match, err := ComparePasswordAndHash(body.Password, user.PasswordHash)
	if err != nil {
		return utils.SendError(ctx, err.Error(), fiber.StatusInternalServerError)
	}

	if !match {
		return utils.SendError(ctx, dto.InvalidPassword, fiber.StatusUnauthorized)
	}

	token, err := generateToken(user.ID, time.Now().Add(time.Hour*24), ctl.jwtSecret)
	if err != nil {
		return utils.SendError(ctx, err.Error(), fiber.StatusInternalServerError)
	}

	resp := &dto.AuthResponse{
		Token: token,
	}

	return ctx.Status(fiber.StatusOK).JSON(&resp)
}
