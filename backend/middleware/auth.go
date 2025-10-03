package middleware

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rahulcodepython/todo-backend/apps/users"
	"github.com/rahulcodepython/todo-backend/backend/response"
)

func Authenticated(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the "Authorization" header from the incoming request.
		authorization := c.Get("Authorization")

		// CRITICAL FIX: Check if the Authorization header is empty FIRST
		if authorization == "" {
			return response.UnauthorizedAccess(c, nil, "Authorization header is missing")
		}

		// Split the authorization header
		authorizationParts := strings.Split(authorization, " ")

		// Check if the header has exactly 2 parts: "Bearer" and "token"
		if len(authorizationParts) != 2 {
			return response.UnauthorizedAccess(c, nil, "Invalid Authorization header format. Expected 'Bearer <token>'")
		}

		// Check if the first part is "Bearer"
		if authorizationParts[0] != "Bearer" {
			return response.UnauthorizedAccess(c, nil, "Authorization type must be 'Bearer'")
		}

		// Extract the token
		token := authorizationParts[1]

		// Check if token is empty (e.g., "Bearer ")
		if token == "" {
			return response.UnauthorizedAccess(c, nil, "Token is missing")
		}

		var count int
		var jwt users.JWT

		// Execute SQL query to find the JWT token
		err := db.QueryRow(
			"SELECT COUNT(*) OVER() AS count, id, token, expires_at FROM jwt_tokens WHERE token = $1",
			token,
		).Scan(&count, &jwt.ID, &jwt.Token, &jwt.ExpiresAt)

		if err != nil {
			return response.InternelServerError(c, err, "Internal Server Error")
		}

		// Check if token exists in database
		if count == 0 {
			return response.UnauthorizedAccess(c, nil, "Invalid token")
		}

		// Check if token has expired
		if jwt.ExpiresAt.Before(time.Now()) {
			// Delete expired token
			_, err := db.Exec(users.DeleteJWTByIdQuery, jwt.ID)
			if err != nil {
				return response.InternelServerError(c, err, "Internal Server Error")
			}
			return response.UnauthorizedAccess(c, nil, "Token has expired. Please login again.")
		}

		// Store JWT in context for use by route handlers
		c.Locals("jwt", jwt)

		// Continue to next handler
		return c.Next()
	}
}
