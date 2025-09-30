package middleware

import (
	// Import the Fiber web framework, which provides the core functionalities for building web applications in Go.
	"github.com/gofiber/fiber/v2"
	// Import the Fiber recover middleware, which is designed to gracefully handle panics that occur during request processing.
	// This prevents the server from crashing and instead returns an appropriate error response.
	"github.com/gofiber/fiber/v2/middleware/recover"
	// Import the application's configuration package. Although not directly used in the current `recover` middleware's options,
	// it's a common practice to pass the configuration to middleware functions for potential future customization
	// (e.g., logging panic details based on environment settings).
	"github.com/rahulcodepython/todo-backend/backend/config"
)

// Recover is a middleware function that configures and returns a Fiber recover instance.
// This middleware catches panics that occur in subsequent handlers within the request-response cycle.
// When a panic is caught, it prevents the server from crashing and typically sends a 500 Internal Server Error response
// to the client, ensuring the application remains stable and available.
// It takes a pointer to the application's configuration struct as a parameter,
// although the current implementation doesn't directly use 'cfg' for recover settings.
func Recover(cfg *config.Config) fiber.Handler {
	// Return a new recover middleware instance with its default configuration.
	// The default behavior is to catch panics and return a 500 status code.
	return recover.New()
}
