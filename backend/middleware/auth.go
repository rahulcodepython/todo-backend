package middleware

import (
	"database/sql" // Import the "database/sql" package to interact with SQL databases.
	"strings"      // Import the "strings" package to perform string manipulation, specifically for splitting the Authorization header.
	"time"         // Import the "time" package to handle time-related operations, such as checking token expiration.

	"github.com/gofiber/fiber/v2"                        // Import the Fiber web framework, which provides the core functionalities for building web applications in Go.
	"github.com/rahulcodepython/todo-backend/apps/users" // Import the "users" package from the application's "apps" directory, specifically to use the `users.JWT` struct.
)

// Authenticated is a Fiber middleware function that checks if an incoming request is authenticated.
// It expects a JWT (JSON Web Token) in the "Authorization" header in the format "Bearer <token>".
// This middleware performs several checks:
// 1. Verifies the presence and format of the Authorization header.
// 2. Checks if the token exists in the database.
// 3. Checks if the token has expired.
// If all checks pass, it stores the JWT information in `c.Locals("jwt")` and allows the request to proceed to the next handler.
// If any check fails, it returns an appropriate HTTP status code and a JSON error response, preventing further processing.
// It takes a database connection (`*sql.DB`) as a parameter to query the `jwt_tokens` table.
func Authenticated(db *sql.DB) fiber.Handler {
	// Return a Fiber handler function that will be executed for each incoming request.
	return func(c *fiber.Ctx) error {
		// Retrieve the "Authorization" header from the incoming request.
		authorization := c.Get("Authorization")
		// Check if the Authorization header is empty.
		if authorization == "" {
			// If it's empty, return an Unauthorized status (401) with a JSON error message.
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,                 // Indicate that the operation was not successful.
				"message": "Unauthorized Access", // Provide a descriptive error message.
			})
		}

		// Split the Authorization header value by space to separate the token type (Bearer) from the actual token.
		authorizationParts := strings.Split(authorization, " ")
		// Extract the token type (e.g., "Bearer").
		header := string(authorizationParts[0])
		// Extract the actual JWT token.
		token := string(authorizationParts[1])

		// Check if the token type is "Bearer". If not, it's an invalid authorization scheme.
		if header != "Bearer" {
			// Return an Unauthorized status (401) with a JSON error message.
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,                 // Indicate that the operation was not successful.
				"message": "Unauthorized Access", // Provide a descriptive error message.
			})
		}

		var count int     // Declare a variable to store the count of matching tokens found in the database.
		var jwt users.JWT // Declare a variable of type `users.JWT` to store the token's details from the database.

		// Execute a SQL query to find a JWT token by its value and retrieve its ID, token string, and expiration time.
		// The `COUNT(*) OVER()` is used to get the total count of rows that would be returned by the query,
		// which helps in checking if a token exists without a separate `SELECT COUNT(*)` query.
		err := db.QueryRow("SELECT COUNT(*) OVER() AS count, id, token, expires_at FROM jwt_tokens WHERE token = $1", token).Scan(&count, &jwt.ID, &jwt.Token, &jwt.ExpiresAt)
		// Check for any database errors during the query execution.
		if err != nil {
			// If an error occurs, return an Internal Server Error status (500) with a JSON error message and the actual error.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,                   // Indicate that the operation was not successful.
				"message": "Internal Server Error", // Provide a generic error message for internal issues.
				"error":   err.Error(),             // Include the specific database error for debugging.
			})
		}

		// After scanning, check if `count` is 0, meaning no matching token was found in the database.
		if count == 0 {
			// If no token is found, return an Unauthorized status (401) with a JSON error message.
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,                 // Indicate that the operation was not successful.
				"message": "Unauthorized Access", // Provide a descriptive error message.
			})
		}

		// Check if the retrieved token's expiration time is before the current time.
		if jwt.ExpiresAt.Before(time.Now()) {
			// If the token has expired, delete it from the `jwt_tokens` table to clean up expired tokens.
			_, err := db.Exec("DELETE FROM jwt_tokens WHERE id = $1", jwt.ID)
			// Check for any database errors during the deletion.
			if err != nil {
				// If an error occurs during deletion, return an Internal Server Error status (500) with a JSON error message.
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,                   // Indicate that the operation was not successful.
					"message": "Internal Server Error", // Provide a generic error message for internal issues.
					"error":   err.Error(),             // Include the specific database error for debugging.
				})
			}

			// After deleting the expired token, return an Unauthorized status (401) with a specific message
			// indicating that the token has expired and the user needs to log in again.
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,                                    // Indicate that the operation was not successful.
				"message": "Token has expired. Please login again.", // Inform the user about the token expiration.
			})
		}

		// If the token is valid and not expired, store the `jwt` struct in Fiber's locals context.
		// This makes the JWT information accessible to subsequent handlers in the request chain.
		c.Locals("jwt", jwt)
		// Call `c.Next()` to pass control to the next middleware or route handler in the chain.
		return c.Next()
	}
}
