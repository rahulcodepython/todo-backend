package middleware

import (
	// Import the Fiber web framework, which provides the core functionalities for building web applications in Go.
	"github.com/gofiber/fiber/v2"
	// Import the Fiber logger middleware, which provides request logging capabilities.
	"github.com/gofiber/fiber/v2/middleware/logger"
	// Import the application's configuration package, although it's not directly used in the current logger configuration,
	// it's often included in middleware functions to allow for future configuration-driven logging behavior.
	"github.com/rahulcodepython/todo-backend/backend/config"
)

// Logger is a middleware function that configures and returns a Fiber logger instance.
// This middleware is responsible for logging details about incoming HTTP requests and their responses,
// which is crucial for debugging, monitoring, and auditing web application activity.
// It takes a pointer to the application's configuration struct as a parameter,
// although the current implementation doesn't directly use 'cfg' for logger settings.
func Logger(cfg *config.Config) fiber.Handler {
	// Return a new logger middleware instance configured with specific options.
	return logger.New(logger.Config{
		// Format specifies the output format for log entries.
		// This string uses placeholders (e.g., ${protocol}, ${ip}, ${status}) that Fiber's logger
		// will replace with actual request and response details for each incoming HTTP request.
		// The "\n" at the end ensures each log entry is on a new line for readability.
		Format: "${protocol}://${ip}:${port} ${status} - ${method} ${path} [${time}]\n",
	})
}
