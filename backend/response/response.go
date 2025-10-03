// This file provides standardized functions for sending API responses.
package response

// "github.com/gofiber/fiber/v2" is a web framework for Go. It is used here to send HTTP responses.
import (
	"github.com/gofiber/fiber/v2"
	// "github.com/rahulcodepython/todo-backend/backend/utils" is a local package that provides the standard response structure.
	"github.com/rahulcodepython/todo-backend/backend/utils"
)

// InternelServerError sends a 500 Internal Server Error response.
// It takes the Fiber context, an error, and a message as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @param err error - The error that occurred.
// @param message string - A message to be included in the response.
// @return error - An error if one occurred while sending the response.
func InternelServerError(c *fiber.Ctx, err error, message string) error {
	// This checks if a custom message is provided.
	if message == "" {
		// If no message is provided, a default message is used.
		message = "Internal Server Error"
	}

	// c.Status() sets the HTTP status code of the response.
	// c.JSON() sends a JSON response.
	return c.Status(fiber.StatusInternalServerError).JSON(utils.Response{
		// Success is set to false to indicate that the request was not successful.
		Success: false,
		// The message is included in the response.
		Message: message,
		// The error message is included in the response.
		Error: err.Error(),
	})
}

// BadInternalResponse sends a 400 Bad Request response.
// It takes the Fiber context, an error, and a message as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @param err error - The error that occurred.
// @param message string - A message to be included in the response.
// @return error - An error if one occurred while sending the response.
func BadInternalResponse(c *fiber.Ctx, err error, message string) error {
	// This checks if a custom message is provided.
	if message == "" {
		// If no message is provided, a default message is used.
		message = "Bad Request"
	}

	// c.Status() sets the HTTP status code of the response.
	// c.JSON() sends a JSON response.
	return c.Status(fiber.StatusBadRequest).JSON(utils.Response{
		// Success is set to false to indicate that the request was not successful.
		Success: false,
		// The message is included in the response.
		Message: message,
		// The error message is included in the response.
		Error: err.Error(),
	})
}

// UnauthorizedAccess sends a 401 Unauthorized response.
// It takes the Fiber context, an error, and a message as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @param err error - The error that occurred.
// @param message string - A message to be included in the response.
// @return error - An error if one occurred while sending the response.
func UnauthorizedAccess(c *fiber.Ctx, err error, message string) error {
	// This checks if a custom message is provided.
	if message == "" {
		// If no message is provided, a default message is used.
		message = "Unauthorized Access"
	}

	// c.Status() sets the HTTP status code of the response.
	// c.JSON() sends a JSON response.
	return c.Status(fiber.StatusUnauthorized).JSON(utils.Response{
		// Success is set to false to indicate that the request was not successful.
		Success: false,
		// The message is included in the response.
		Message: message,
		// The error message is included in the response.
		Error: err.Error(),
	})
}

// NotFound sends a 404 Not Found response.
// It takes the Fiber context, an error, and a message as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @param err error - The error that occurred.
// @param message string - A message to be included in the response.
// @return error - An error if one occurred while sending the response.
func NotFound(c *fiber.Ctx, err error, message string) error {
	// This checks if a custom message is provided.
	if message == "" {
		// If no message is provided, a default message is used.
		message = "Not Found"
	}

	// c.Status() sets the HTTP status code of the response.
	// c.JSON() sends a JSON response.
	return c.Status(fiber.StatusNotFound).JSON(utils.Response{
		// Success is set to false to indicate that the request was not successful.
		Success: false,
		// The message is included in the response.
		Message: message,
		// The error message is included in the response.
		Error: err.Error(),
	})
}

// BadResponse sends a 400 Bad Request response.
// It takes the Fiber context and a message as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @param message string - A message to be included in the response.
// @return error - An error if one occurred while sending the response.
func BadResponse(c *fiber.Ctx, message string) error {
	// This checks if a custom message is provided.
	if message == "" {
		// If no message is provided, a default message is used.
		message = "Bad Request"
	}

	// c.Status() sets the HTTP status code of the response.
	// c.JSON() sends a JSON response.
	return c.Status(fiber.StatusBadRequest).JSON(utils.Response{
		// Success is set to false to indicate that the request was not successful.
		Success: false,
		// The message is included in the response.
		Message: message,
	})
}

// OKResponse sends a 200 OK response.
// It takes the Fiber context, a message, and data as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @param message string - A message to be included in the response.
// @param data interface{} - The data to be included in the response.
// @return error - An error if one occurred while sending the response.
func OKResponse(c *fiber.Ctx, message string, data interface{}) error {
	// c.Status() sets the HTTP status code of the response.
	// c.JSON() sends a JSON response.
	return c.Status(fiber.StatusOK).JSON(utils.Response{
		// Success is set to true to indicate that the request was successful.
		Success: true,
		// The message is included in the response.
		Message: message,
		// The data is included in the response.
		Data: data,
	})
}

// OKCreatedResponse sends a 201 Created response.
// It takes the Fiber context, a message, and data as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @param message string - A message to be included in the response.
// @param data interface{} - The data to be included in the response.
// @return error - An error if one occurred while sending the response.
func OKCreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	// c.Status() sets the HTTP status code of the response.
	// c.JSON() sends a JSON response.
	return c.Status(fiber.StatusCreated).JSON(utils.Response{
		// Success is set to true to indicate that the request was successful.
		Success: true,
		// The message is included in the response.
		Message: message,
		// The data is included in the response.
		Data: data,
	})
}

// TooManyRequests sends a 429 Too Many Requests response.
// It takes the Fiber context and a message as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @param message string - A message to be included in the response.
// @return error - An error if one occurred while sending the response.
func TooManyRequests(c *fiber.Ctx, message string) error {
	// c.Status() sets the HTTP status code of the response.
	// c.JSON() sends a JSON response.
	return c.Status(fiber.StatusTooManyRequests).JSON(utils.Response{
		// Success is set to false to indicate that the request was not successful.
		Success: false,
		// The message is included in the response.
		Message: message,
	})
}