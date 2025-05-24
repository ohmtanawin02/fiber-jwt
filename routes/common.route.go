package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ohmtanawin02/go-fiber-jwt-v2/configs"
	"github.com/ohmtanawin02/go-fiber-jwt-v2/handlers"
	"github.com/ohmtanawin02/go-fiber-jwt-v2/middlewares"
)

func SetupRoutes(app *fiber.App, handler *handlers.Handler) {
	authMiddleware := middlewares.NewAuthMiddleware(configs.Configs.SECRET_KEY)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	RegisterUserRoutes(v1, handler, authMiddleware)
	RegisterAuthPublicRoutes(v1, handler, authMiddleware)
}
