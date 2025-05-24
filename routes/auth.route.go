package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ohmtanawin02/go-fiber-jwt-v2/handlers"
	"github.com/ohmtanawin02/go-fiber-jwt-v2/middlewares"
)

func RegisterAuthPublicRoutes(router fiber.Router, handler *handlers.Handler, authMiddleware *middlewares.AuthMiddleware) {
	authPublic := router.Group("/public")

	authPublic.Post("/register", handler.RegisterUser)
}
