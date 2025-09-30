package middleware

import (
	// Import the Fiber web framework, which provides the core functionalities for building web applications in Go.
	"github.com/gofiber/fiber/v2"
	// Import the Fiber CORS middleware, which handles Cross-Origin Resource Sharing.
	"github.com/gofiber/fiber/v2/middleware/cors"
	// Import the application's configuration package to access CORS-related settings.
	"github.com/rahulcodepython/todo-backend/backend/config"
)

// Cors is a middleware function that configures Cross-Origin Resource Sharing (CORS) for the Fiber application.
// CORS is a security feature implemented by web browsers that restricts web pages from making requests
// to a different domain than the one that served the web page. This middleware allows the server to specify
// which origins are permitted to access its resources.
// It takes a pointer to the application's configuration struct as a parameter to retrieve CORS settings.
func Cors(cfg *config.Config) fiber.Handler {
	// Return a new CORS middleware instance configured with the specified options.
	return cors.New(cors.Config{
		// AllowOrigins specifies a comma-separated list of origins that are allowed to make cross-origin requests.
		// This value is retrieved from the application's configuration (cfg.CORS.CorsOrigins).
		AllowOrigins: cfg.CORS.CorsOrigins,
		// AllowHeaders specifies a comma-separated list of HTTP headers that can be used when making the actual request.
		AllowHeaders: "Origin, Content-Type, Accept",
		// Next is a function that defines whether the middleware should be skipped for the current request.
		// In this configuration, the CORS middleware is skipped if the client's IP address matches the server's host IP.
		// This is typically used to bypass CORS checks for requests originating from the same server,
		// which might be useful in certain deployment scenarios or for internal health checks.
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == cfg.Server.Host
		},
	})
}
