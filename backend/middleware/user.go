// This file defines middleware for handling authenticated users.
package middleware

// "database/sql" provides a generic SQL interface. It is used here to query the database.
import (
	"database/sql"

	// "github.com/gofiber/fiber/v2" is a web framework for Go. It is used here to create middleware.
	"github.com/gofiber/fiber/v2"
	// "github.com/rahulcodepython/todo-backend/apps/users" is a local package that contains user-related models and queries.
	"github.com/rahulcodepython/todo-backend/apps/users"
	// "github.com/rahulcodepython/todo-backend/backend/response" is a local package that provides standardized API responses.
	"github.com/rahulcodepython/todo-backend/backend/response"
)

// AuthenticatedUser is a middleware that retrieves the authenticated user's data from the database.
// It should be used after the Authenticated middleware.
// It takes a database connection as input and returns a Fiber handler.
//
// @param db *sql.DB - The database connection.
// @return fiber.Handler - The Fiber handler.
func AuthenticatedUser(db *sql.DB) fiber.Handler {
	// This returns a new Fiber handler.
	return func(c *fiber.Ctx) error {
		// jwtInterface is the JWT object retrieved from the local context.
		jwtInterface := c.Locals("jwt")

		// This checks if the JWT exists in the context.
		if jwtInterface == nil {
			// If the JWT does not exist, it returns an unauthorized access response.
			return response.UnauthorizedAccess(c, nil, "Authentication required")
		}

		// jwt is the JWT object after type assertion.
		jwt, ok := jwtInterface.(users.JWT)
		// This checks if the type assertion was successful.
		if !ok {
			// If the type assertion fails, it returns an internal server error response.
			return response.InternelServerError(c, nil, "Invalid authentication data")
		}

		// user is a variable that will hold the user's data.
		var user users.User

		// err is the result of querying the database for the user's profile.
		// db.QueryRow() executes a query that is expected to return at most one row.
		err := db.QueryRow(
			// users.GetUserProfileByJWTQuery is the SQL query to retrieve the user's profile.
			users.GetUserProfileByJWTQuery,
			// jwt.ID is the ID of the JWT.
			jwt.ID,
		).Scan(
			// The following are the fields to be scanned from the database row.
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Image,
			&user.Password,
			&user.JWT,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		// This checks if an error occurred while querying the database.
		if err != nil {
			// This checks if the error is sql.ErrNoRows.
			if err == sql.ErrNoRows {
				// If no user is found, it returns a not found response.
				return response.NotFound(c, err, "User not found")
			}
			// For any other error, it returns an internal server error response.
			return response.InternelServerError(c, err, "Error fetching user data")
		}

		// The user's data is stored in the local context.
		c.Locals("user", user)

		// c.Next() calls the next middleware in the chain.
		return c.Next()
	}
}