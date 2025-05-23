package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oklog/ulid/v2"
	"orchestrator/internal/controllers/dto"
	"orchestrator/internal/db/users"
	"orchestrator/internal/utils"
	"time"
)

// register godoc
// @Summary      Зарегистрировать новый аккаунт
// @Description  Создать новый аккаунт с помощью логина и пароля. Возвращает JWT API Token при успешном запросе
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body  dto.AuthRequest true  "Данные для регистрации"
// @Success      201  {object}  dto.AuthResponse
// @Failure      400  {object}  dto.ApiError
// @Failure      409  {object}  dto.ApiError
// @Failure      422  {object}  dto.ApiError
// @Failure      500  {object}  dto.ApiError
// @Router       /register [post]
func (ctl *Controller) register(ctx fiber.Ctx) error {
	var body dto.AuthRequest
	if err := ctx.Bind().JSON(&body); err != nil {
		return utils.SendError(ctx, dto.InvalidData, fiber.StatusUnprocessableEntity)
	}

	if len(body.Login) < 1 || len(body.Password) < 1 {
		return utils.SendError(ctx, dto.InvalidData, fiber.StatusBadRequest)
	}

	_, err := ctl.usersRepo.GetByName(ctx.Context(), body.Login)
	if err == nil {
		return utils.SendError(ctx, dto.Conflict, fiber.StatusConflict)
	}

	passwordHash, err := CreateHash(body.Password, DefaultParams)
	if err != nil {
		return utils.SendError(ctx, err.Error(), fiber.StatusInternalServerError)
	}

	userID := ulid.Make().String()

	user := users.CreateParams{
		ID:           userID,
		Username:     body.Login,
		PasswordHash: passwordHash,
	}
	if err := ctl.usersRepo.Create(ctx.Context(), user); err != nil {
		return utils.SendError(ctx, err.Error(), fiber.StatusInternalServerError)
	}

	token, err := generateToken(userID, time.Now().Add(time.Hour*24), ctl.jwtSecret)
	if err != nil {
		return utils.SendError(ctx, err.Error(), fiber.StatusInternalServerError)
	}

	resp := &dto.AuthResponse{
		Token: token,
	}

	return ctx.Status(fiber.StatusCreated).JSON(&resp)
}
