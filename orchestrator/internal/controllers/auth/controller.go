package auth

import (
	"github.com/gofiber/fiber/v3"
	"orchestrator/internal/db/users"
)

type Controller struct {
	usersRepo *users.Queries
	jwtSecret string
}

func New(usersRepo *users.Queries, jwtSecret string) *Controller {
	return &Controller{usersRepo: usersRepo, jwtSecret: jwtSecret}
}

func (ctl *Controller) Setup(group fiber.Router) {
	group.Post("/login", ctl.login)
	group.Post("/register", ctl.register)
}
