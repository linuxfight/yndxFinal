package controllers

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"orchestrator/internal/controllers/auth"
	"orchestrator/internal/controllers/expr"
	jwtware "orchestrator/internal/controllers/middlewares/jwt"
	"orchestrator/internal/controllers/middlewares/recoverer"
	"orchestrator/internal/controllers/middlewares/swagger"
	"orchestrator/internal/db"
	"orchestrator/internal/db/expressions"
	"orchestrator/internal/db/users"
	"orchestrator/internal/utils"
	"os"
)

const testingSecret = "not_v3ry_s3cR3T"

func NewFiber(userRepo *users.Queries, exprRepo *expressions.Queries, cache *db.Cache) *fiber.App {
	// create fiber app
	cfg := fiber.Config{
		JSONDecoder: sonic.Unmarshal,
		JSONEncoder: sonic.Marshal,
	}
	app := fiber.New(cfg)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = testingSecret
	}

	// set up middlewares
	app.Get(healthcheck.DefaultStartupEndpoint, healthcheck.NewHealthChecker())
	app.Use(cors.New())
	app.Use(recoverer.New())
	app.Use(swagger.New(swagger.Config{
		FilePath: "./docs/swagger.json",
		Path:     "./docs",
		Title:    "Swagger API Docs",
	}))

	// create router
	group := app.Group("/api/v1")

	// set up auth
	authCtl := auth.New(userRepo, jwtSecret)
	authCtl.Setup(group)

	// set up expr
	authWare := jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			return utils.SendError(ctx, err.Error(), fiber.StatusUnauthorized)
		},
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256, // HMAC256 signing method
			Key:    []byte(jwtSecret),
		},
	})
	solverCtl := expr.New(exprRepo, cache)
	solverCtl.Setup(group, authWare)

	return app
}
