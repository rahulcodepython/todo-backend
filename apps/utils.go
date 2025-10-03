package apps

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rahulcodepython/todo-backend/backend/response"
)

// BodyParser is a generic helper function to parse the request body into a given struct.
// Using generics [T any] allows this function to work with any struct type (e.g., registerUserRequest, loginUserRequest).
func BodyParser[T any](c *fiber.Ctx, body *T) error {
	// c.BodyParser attempts to parse the incoming request body and unmarshal it into the 'body' struct.
	if err := c.BodyParser(body); err != nil {
		// If parsing fails, it sends a standardized bad request response to the client.
		// This centralizes the parsing and error handling logic.
		return response.BadInternalResponse(c, err, "Invalid request body")
	}
	// If parsing is successful, return nil to indicate no error occurred.
	return nil
}
