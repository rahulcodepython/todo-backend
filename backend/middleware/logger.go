// This file defines a middleware for logging HTTP requests.
package middleware

// "github.com/gofiber/fiber/v2" is a web framework for Go. It is used here to create middleware.
import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/logger" is a middleware that logs requests.
	"github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/rahulcodepython/todo-backend/backend/config" is a local package that provides access to the application configuration.
	"github.com/rahulcodepython/todo-backend/backend/config"
)

// Logger is a middleware that logs HTTP requests.
// It takes the application configuration as input and returns a Fiber handler.
//
// @param cfg *config.Config - The application configuration.
// @return fiber.Handler - The Fiber handler.
func Logger(cfg *config.Config) fiber.Handler {
	// logger.New() returns a new logger middleware with the specified configuration.
	return logger.New(logger.Config{
		// Format is the format of the log message.
		Format: "[${time}] ${protocol}://${ip}:${port} - ${method} : ${status} | ${path} | ${latency} \n", // Time is the timestamp of the log entry.
		// Protocol is the protocol used for the request (e.g., HTTP/1.1).
		// IP is the IP address of the client.
		// Port is the port number of the server.
		// Method is the HTTP method of the request (e.g., GET, POST).
		// Status is the HTTP status code of the response.
		// Path is the URL path of the request.
		// Latency is the time taken to process the request.

	})
}
