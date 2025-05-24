package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ohmtanawin02/go-fiber-jwt-v2/handlers"
	"github.com/ohmtanawin02/go-fiber-jwt-v2/middlewares"
)

func RegisterUserRoutes(router fiber.Router, handler *handlers.Handler, authMiddleware *middlewares.AuthMiddleware) {
	users := router.Group("/users")
	protected := users.Group("/")
	// protected.Use(authMiddleware.AuthenticateJWT())

	protected.Get("/:id", handler.GetUser)
	protected.Get("/", handler.GetUsers)
	protected.Delete("/:id", handler.SoftDeleteUser)

}
