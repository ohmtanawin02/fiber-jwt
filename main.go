package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ohmtanawin02/go-fiber-jwt-v2/configs"
	"github.com/ohmtanawin02/go-fiber-jwt-v2/handlers"
	"github.com/ohmtanawin02/go-fiber-jwt-v2/middlewares"
	"github.com/ohmtanawin02/go-fiber-jwt-v2/routes"
)

func main() {
	config := configs.LoadConfigs()
	config.ConnectDatabase()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	middlewares.SetupMiddleware(app)

	handler := handlers.NewHandler(config.MongoDB, config.MongoDB_Database)

	routes.SetupRoutes(app, handler)

	// Start server
	port := fmt.Sprintf(":%d", config.Port)
	log.Printf("Server starting on port %s", port)

	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
