package expr

import (
	"github.com/gofiber/fiber/v3"
	"orchestrator/internal/db"
	"orchestrator/internal/db/expressions"
)

type Controller struct {
	exprRepo *expressions.Queries
	tasks    *db.Cache
}

func New(exprs *expressions.Queries, cache *db.Cache) *Controller {
	return &Controller{exprRepo: exprs, tasks: cache}
}

func (ctl *Controller) Setup(group fiber.Router, authWare fiber.Handler) {
	group.Get("/expressions", ctl.list, authWare)
	group.Get("/expressions/:id", ctl.GetById, authWare)
	group.Post("/calculate", ctl.calculate, authWare)
}
