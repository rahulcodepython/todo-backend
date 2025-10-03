// This file defines a middleware for handling authentication.
package middleware

// "database/sql" provides a generic SQL interface. It is used here to query the database.
import (
	"database/sql"
	// "strings" provides functions for working with strings. It is used here to split the Authorization header.
	"strings"
	// "time" provides functions for working with time. It is used here to check if a JWT has expired.
	"time"

	// "github.com/gofiber/fiber/v2" is a web framework for Go. It is used here to create middleware.
	"github.com/gofiber/fiber/v2"
	// "github.com/rahulcodepython/todo-backend/apps/users" is a local package that contains user-related models and queries.
	"github.com/rahulcodepython/todo-backend/apps/users"
	// "github.com/rahulcodepython/todo-backend/backend/response" is a local package that provides standardized API responses.
	"github.com/rahulcodepython/todo-backend/backend/response"
)

// Authenticated is a middleware that checks if a user is authenticated.
// It takes a database connection as input and returns a Fiber handler.
//
// @param db *sql.DB - The database connection.
// @return fiber.Handler - The Fiber handler.
func Authenticated(db *sql.DB) fiber.Handler {
	// This returns a new Fiber handler.
	return func(c *fiber.Ctx) error {
		// authorization is the value of the "Authorization" header.
		authorization := c.Get("Authorization")

		// This checks if the Authorization header is empty.
		if authorization == "" {
			// If the header is empty, it returns an unauthorized access response.
			return response.UnauthorizedAccess(c, nil, "Authorization header is missing")
		}

		// authorizationParts is a slice of strings that contains the parts of the Authorization header.
		authorizationParts := strings.Split(authorization, " ")

		// This checks if the header has exactly two parts.
		if len(authorizationParts) != 2 {
			// If the header does not have two parts, it returns an unauthorized access response.
			return response.UnauthorizedAccess(c, nil, "Invalid Authorization header format. Expected 'Bearer <token>'")
		}

		// This checks if the first part of the header is "Bearer".
		if authorizationParts[0] != "Bearer" {
			// If the first part is not "Bearer", it returns an unauthorized access response.
			return response.UnauthorizedAccess(c, nil, "Authorization type must be 'Bearer'")
		}

		// token is the second part of the Authorization header.
		token := authorizationParts[1]

		// This checks if the token is empty.
		if token == "" {
			// If the token is empty, it returns an unauthorized access response.
			return response.UnauthorizedAccess(c, nil, "Token is missing")
		}

		// count is a variable that will hold the number of rows returned by the query.
		var count int
		// jwt is a variable that will hold the JWT data.
		var jwt users.JWT

		// err is the result of querying the database for the JWT.
		// db.QueryRow() executes a query that is expected to return at most one row.
		err := db.QueryRow(
			// This is the SQL query to retrieve the JWT.
			"SELECT COUNT(*) OVER() AS count, id, token, expires_at FROM jwt_tokens WHERE token = $1",
			// token is the token from the Authorization header.
			token,
		).Scan(&count, &jwt.ID, &jwt.Token, &jwt.ExpiresAt)

		// This checks if an error occurred while querying the database.
		if err != nil {
			// If an error occurs, it returns an internal server error response.
			return response.InternelServerError(c, err, "Internal Server Error")
		}

		// This checks if the token exists in the database.
		if count == 0 {
			// If the token does not exist, it returns an unauthorized access response.
			return response.UnauthorizedAccess(c, nil, "Invalid token")
		}

		// This checks if the token has expired.
		if jwt.ExpiresAt.Before(time.Now()) {
			// If the token has expired, it is deleted from the database.
			_, err := db.Exec(users.DeleteJWTByIdQuery, jwt.ID)
			// This checks if an error occurred while deleting the token.
			if err != nil {
				// If an error occurs, it returns an internal server error response.
				return response.InternelServerError(c, err, "Internal Server Error")
			}
			// It then returns an unauthorized access response.
			return response.UnauthorizedAccess(c, nil, "Token has expired. Please login again.")
		}

		// The JWT data is stored in the local context.
		c.Locals("jwt", jwt)

		// c.Next() calls the next middleware in the chain.
		return c.Next()
	}
}