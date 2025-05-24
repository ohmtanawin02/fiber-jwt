package handlers

import (
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ohmtanawin02/go-fiber-jwt-v2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *Handler) GetUsers(c *fiber.Ctx) error {
	userCollection := h.GetCollection("users")
	request := models.GetUserPaginateRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
		})
	}
	searchQuery := request.Search
	emailFilter := request.Email

	page := max(request.Page, 1)
	limit := max(request.Limit, 10)

	filter := bson.M{
		"is_deleted": false,
	}

	if searchQuery != "" {
		searchRegex := primitive.Regex{Pattern: regexp.QuoteMeta(searchQuery), Options: "i"}
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": searchRegex}},
			{"email": bson.M{"$regex": searchRegex}},
			{"username": bson.M{"$regex": searchRegex}},
		}
	}

	if emailFilter != "" {
		filter["email"] = emailFilter
	}

	total, err := userCollection.CountDocuments(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to count users",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	skip := (page - 1) * limit

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"created_at": -1})

	cursor, err := userCollection.Find(c.Context(), filter, findOptions)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to fetch users",
		})
	}

	var users []models.User
	if err := cursor.All(c.Context(), &users); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to decode users",
		})
	}

	userResult := make([]fiber.Map, 0)

	for _, user := range users {
		userResult = append(userResult, fiber.Map{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.PaginationResponse{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Data:       userResult,
	})
}

func (h *Handler) GetUser(c *fiber.Ctx) error {
	userCollection := h.GetCollection("users")
	request := models.GetUserRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
		})
	}
	objectID, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID format",
		})
	}

	var user models.User
	err = userCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
		IsDeleted: user.IsDeleted,
	})
}

func (h *Handler) HardDeleteUser(c *fiber.Ctx) error {
	userCollection := h.GetCollection("users")
	id := c.Params("id")
	fmt.Printf("User ID: %s\n", id)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID format",
		})
	}

	fmt.Printf("User objectID: %s\n", objectID)

	_, err = userCollection.DeleteOne(c.Context(), bson.M{"_id": objectID})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Delete success",
	})
}

func (h *Handler) SoftDeleteUser(c *fiber.Ctx) error {
	userCollection := h.GetCollection("users")
	fmt.Printf("SOFT DELETE")
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID format",
		})
	}

	update := bson.M{
		"$set": bson.M{
			"deletedAt": time.Now(),
		},
	}

	filter := bson.M{
		"_id":       objectID,
		"deletedAt": bson.M{"$exists": nil},
	}

	result, err := userCollection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to delete user",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found or already deleted",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User soft deleted successfully",
	})
}
