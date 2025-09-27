package users

import (
	"database/sql"
	"strings"
	"time"

	"github.com/rahulcodepython/todo-backend/backend/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": utils.ErrMissingAuthHeader,
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": utils.ErrInvalidAuthHeader,
			})
		}

		tokenString := parts[1]

		var userID string
		var expiresAt time.Time

		err := db.QueryRow(GetJWTByTokenQuery, tokenString).Scan(&userID, &expiresAt)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": utils.ErrInvalidToken,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Database error",
			})
		}

		if time.Now().After(expiresAt) {
			// Optional: Delete expired token from DB
			db.Exec(DeleteJWTByTokenQuery, tokenString)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": utils.ErrTokenExpired,
			})
		}

		var user User
		err = db.QueryRow(GetUserByIDQuery, userID).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": utils.ErrUserNotFound,
			})
		}

		// Store user in context for downstream handlers
		c.Locals(utils.CtxUserKey, &user)

		return c.Next()
	}
}
