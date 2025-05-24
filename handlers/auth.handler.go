package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ohmtanawin02/go-fiber-jwt-v2/models"
	"github.com/ohmtanawin02/go-fiber-jwt-v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) RegisterUser(c *fiber.Ctx) error {
	userCollection := h.GetCollection("users")
	user := models.User{}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
		})
	}

	if user.Email == "" || user.Name == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "email, name, password is required",
		})
	}

	var existingUser models.User
	err := userCollection.FindOne(c.Context(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.Password = hashedPassword
	user.IsDeleted = false

	_, err = userCollection.InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}
