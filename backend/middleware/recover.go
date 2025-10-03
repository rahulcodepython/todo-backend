// This file defines a middleware for recovering from panics.
package middleware

// "github.com/gofiber/fiber/v2" is a web framework for Go. It is used here to create middleware.
import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/recover" is a middleware that recovers from panics anywhere in the stack chain.
	"github.com/gofiber/fiber/v2/middleware/recover"
	// "github.com/rahulcodepython/todo-backend/backend/config" is a local package that provides access to the application configuration.
	"github.com/rahulcodepython/todo-backend/backend/config"
)

// Recover is a middleware that recovers from panics.
// It takes the application configuration as input and returns a Fiber handler.
//
// @param cfg *config.Config - The application configuration.
// @return fiber.Handler - The Fiber handler.
func Recover(cfg *config.Config) fiber.Handler {
	// recover.New() returns a new recover middleware with default configuration.
	return recover.New()
}