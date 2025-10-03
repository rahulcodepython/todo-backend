// middleware/user.go
package middleware

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/rahulcodepython/todo-backend/apps/users"
	"github.com/rahulcodepython/todo-backend/backend/response"
)

// AuthenticatedUser is a Fiber middleware function that loads the authenticated user's data.
// This middleware should be used AFTER the Authenticated middleware.
// It expects a JWT object to be stored in c.Locals("jwt") by the Authenticated middleware.
// This middleware performs the following:
// 1. Retrieves the JWT from context (placed by Authenticated middleware).
// 2. Queries the database to fetch the user's profile using the JWT ID.
// 3. Stores the user data in c.Locals("user") for use by subsequent handlers.
// If any check fails, it returns an appropriate HTTP status code and JSON error response.
func AuthenticatedUser(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the JWT from the local context
		// This should have been set by the Authenticated middleware
		jwtInterface := c.Locals("jwt")

		// Check if JWT exists in context
		if jwtInterface == nil {
			return response.UnauthorizedAccess(c, nil, "Authentication required")
		}

		// Type assert the interface to users.JWT struct
		jwt, ok := jwtInterface.(users.JWT)
		if !ok {
			return response.InternelServerError(c, nil, "Invalid authentication data")
		}

		// Declare a variable to hold the user data
		var user users.User

		// Query the database to fetch user profile based on JWT ID
		// This joins the users table with jwt_tokens table to get the user associated with this JWT
		err := db.QueryRow(
			users.GetUserProfileByJWTQuery,
			jwt.ID,
		).Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Image,
			&user.Password,
			&user.JWT,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		// Handle database errors
		if err != nil {
			if err == sql.ErrNoRows {
				// If no user is found, it means the JWT is valid but the user doesn't exist
				// This could happen if the user was deleted after the JWT was issued
				return response.NotFound(c, err, "User not found")
			}
			// For other database errors, return internal server error
			return response.InternelServerError(c, err, "Error fetching user data")
		}

		// Store the user data in the local context for access by route handlers
		c.Locals("user", user)

		// Continue to the next handler
		return c.Next()
	}
}
