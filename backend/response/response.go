package response

import (
	"github.com/gofiber/fiber/v2" // Import the Fiber web framework, which provides the core functionalities for building web applications in Go.
	"github.com/rahulcodepython/todo-backend/backend/utils"
	// Import the application's utility package, specifically for the `Response` struct.
)

// InternelServerError constructs and sends a 500 Internal Server Error response.
// This function is used when an unexpected error occurs on the server side that prevents
// the request from being fulfilled. It ensures a consistent error response format.
//
// Parameters:
// - c: The Fiber context, providing access to the HTTP request and response.
// - err: The actual error object that occurred, used for detailed logging or debugging.
// - message: A human-readable message describing the error. If empty, a default message is used.
//
// Returns:
// - An error, typically from `c.JSON`, indicating the result of sending the response.
func InternelServerError(c *fiber.Ctx, err error, message string) error {
	// Check if a custom message is provided.
	if message == "" { // If the message is empty,
		message = "Internal Server Error" // set a default "Internal Server Error" message.
	}

	return c.Status(fiber.StatusInternalServerError).JSON(utils.Response{ // Set the HTTP status code to 500 (Internal Server Error) and send a JSON response.
		Success: false,       // Indicate that the operation was not successful.
		Message: message,     // Include the descriptive message (either custom or default).
		Error:   err.Error(), // Include the detailed error message from the `err` object.
	})
}

// BadInternalResponse constructs and sends a 400 Bad Request response, typically for internal validation failures.
// This function is used when the server detects an issue with the request that is not directly
// related to user input format but rather an internal processing error that should result in a bad request.
//
// Parameters:
// - c: The Fiber context.
// - err: The error object that occurred, providing details about the bad request.
// - message: A human-readable message. If empty, a default "Bad Request" message is used.
//
// Returns:
// - An error from `c.JSON`.
func BadInternalResponse(c *fiber.Ctx, err error, message string) error {
	// Check if a custom message is provided.
	if message == "" { // If the message is empty,
		message = "Bad Request" // set a default "Bad Request" message.
	}

	return c.Status(fiber.StatusBadRequest).JSON(utils.Response{ // Set the HTTP status code to 400 (Bad Request) and send a JSON response.
		Success: false,       // Indicate that the operation was not successful.
		Message: message,     // Include the descriptive message (either custom or default).
		Error:   err.Error(), // Include the detailed error message from the `err` object.
	})
}

// UnauthorizedAccess constructs and sends a 401 Unauthorized response.
// This function is used when the client attempts to access a protected resource
// without valid authentication credentials or with expired/invalid credentials.
//
// Parameters:
// - c: The Fiber context.
// - err: The error object, providing details about why the access was unauthorized.
// - message: A human-readable message. If empty, a default "Unauthorized Access" message is used.
//
// Returns:
// - An error from `c.JSON`.
func UnauthorizedAccess(c *fiber.Ctx, err error, message string) error {
	// Check if a custom message is provided.
	if message == "" { // If the message is empty,
		message = "Unauthorized Access" // set a default "Unauthorized Access" message.
	}

	return c.Status(fiber.StatusUnauthorized).JSON(utils.Response{ // Set the HTTP status code to 401 (Unauthorized) and send a JSON response.
		Success: false,       // Indicate that the operation was not successful.
		Message: message,     // Include the descriptive message (either custom or default).
		Error:   err.Error(), // Include the detailed error message from the `err` object.
	})
}

// NotFound constructs and sends a 404 Not Found response.
// This function is used when the requested resource could not be found on the server.
//
// Parameters:
// - c: The Fiber context.
// - err: The error object, providing details about why the resource was not found.
// - message: A human-readable message. If empty, a default "Not Found" message is used.
//
// Returns:
// - An error from `c.JSON`.
func NotFound(c *fiber.Ctx, err error, message string) error {
	// Check if a custom message is provided.
	if message == "" { // If the message is empty,
		message = "Not Found" // set a default "Not Found" message.
	}

	return c.Status(fiber.StatusNotFound).JSON(utils.Response{ // Set the HTTP status code to 404 (Not Found) and send a JSON response.
		Success: false,       // Indicate that the operation was not successful.
		Message: message,     // Include the descriptive message (either custom or default).
		Error:   err.Error(), // Include the detailed error message from the `err` object.
	})
}

// BadResponse constructs and sends a 400 Bad Request response for client-side input errors.
// This function is used when the client sends a request that the server cannot process
// due to invalid syntax, missing parameters, or other client-side issues.
//
// Parameters:
// - c: The Fiber context.
// - message: A human-readable message explaining the bad request. If empty, a default message is used.
//
// Returns:
// - An error from `c.JSON`.
func BadResponse(c *fiber.Ctx, message string) error {
	// Check if a custom message is provided.
	if message == "" { // If the message is empty,
		message = "Bad Request" // set a default "Bad Request" message.
	}

	return c.Status(fiber.StatusBadRequest).JSON(utils.Response{ // Set the HTTP status code to 400 (Bad Request) and send a JSON response.
		Success: false,   // Indicate that the operation was not successful.
		Message: message, // Include the descriptive message (either custom or default).
	})
}

// OKResponse constructs and sends a 200 OK response for successful operations.
// This function is used when a request has been successfully processed and the server
// is returning the requested data or a success confirmation.
//
// Parameters:
// - c: The Fiber context.
// - message: A human-readable message confirming the success of the operation.
// - data: The actual payload to be returned to the client. Can be any type.
//
// Returns:
// - An error from `c.JSON`.
func OKResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(utils.Response{ // Set the HTTP status code to 200 (OK) and send a JSON response.
		Success: true,    // Indicate that the operation was successful.
		Message: message, // Include the success message.
		Data:    data,    // Include the actual data payload.
	})
}

func TooManyRequests(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(utils.Response{
		Success: false,
		Message: message,
	})
}
