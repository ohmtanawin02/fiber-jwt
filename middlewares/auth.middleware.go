package middlewares

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt"
)

type AuthMiddleware struct {
	accessSecret string
}

func NewAuthMiddleware(accessSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		accessSecret: accessSecret,
	}
}

func (m *AuthMiddleware) AuthenticateJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "missing authorization header",
				"error":   "unauthorized",
			})
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid authorization header format",
				"error":   "unauthorized",
			})
		}

		tokenString := headerParts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.accessSecret), nil
		})

		if err != nil {
			log.Error("Error parsing token:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid token",
				"error":   "unauthorized",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid token claims",
				"error":   "unauthorized",
			})
		}

		if exp, ok := claims["exp"].(float64); ok {
			if float64(time.Now().Unix()) > exp {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "token expired",
					"error":   "unauthorized",
				})
			}
		}

		if userID, ok := claims["id"]; ok {
			c.Locals("user_id", userID)
		}

		if username, ok := claims["username"]; ok {
			c.Locals("username", username)
		}

		if email, ok := claims["email"]; ok {
			c.Locals("email", email)
		}

		return c.Next()
	}
}

func GetUserIDFromContext(c *fiber.Ctx) (string, bool) {
	userID := c.Locals("user_id")
	if userID == nil {
		return "", false
	}

	if id, ok := userID.(string); ok {
		return id, true
	}

	return "", false
}
