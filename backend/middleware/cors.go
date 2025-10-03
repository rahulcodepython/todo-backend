// This file defines a middleware for handling Cross-Origin Resource Sharing (CORS).
package middleware

// "github.com/gofiber/fiber/v2" is a web framework for Go. It is used here to create middleware.
import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/cors" is a middleware that provides CORS functionality.
	"github.com/gofiber/fiber/v2/middleware/cors"
	// "github.com/rahulcodepython/todo-backend/backend/config" is a local package that provides access to the application configuration.
	"github.com/rahulcodepython/todo-backend/backend/config"
)

// Cors is a middleware that handles CORS.
// It takes the application configuration as input and returns a Fiber handler.
//
// @param cfg *config.Config - The application configuration.
// @return fiber.Handler - The Fiber handler.
func Cors(cfg *config.Config) fiber.Handler {
	// cors.New() returns a new CORS middleware with the specified configuration.
	return cors.New(cors.Config{
		// AllowOrigins is a list of origins that are allowed to make cross-origin requests.
		AllowOrigins: cfg.CORS.CorsOrigins,
		// AllowHeaders is a list of headers that are allowed in cross-origin requests.
		AllowHeaders: "Origin, Content-Type, Accept",
		// Next is a function that determines whether to skip this middleware.
		Next: func(c *fiber.Ctx) bool {
			// The middleware is skipped if the request is coming from the server itself.
			return c.IP() == cfg.Server.Host
		},
	})
}