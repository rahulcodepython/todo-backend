// This file defines middleware for rate limiting.
package middleware

// "time" provides functions for working with time. It is used here to set the expiration time for the rate limiter.
import (
	"time"

	// "github.com/gofiber/fiber/v2" is a web framework for Go. It is used here to create middleware.
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/limiter" is a middleware that provides rate limiting.
	"github.com/gofiber/fiber/v2/middleware/limiter"
	// "github.com/rahulcodepython/todo-backend/backend/config" is a local package that provides access to the application configuration.
	"github.com/rahulcodepython/todo-backend/backend/config"
	// "github.com/rahulcodepython/todo-backend/backend/response" is a local package that provides standardized API responses.
	"github.com/rahulcodepython/todo-backend/backend/response"
)

// GeneralAPILimiter is a middleware that provides general rate limiting for the API.
// It takes the application configuration as input and returns a Fiber handler.
//
// @param cfg *config.Config - The application configuration.
// @return fiber.Handler - The Fiber handler.
func GeneralAPILimiter(cfg *config.Config) fiber.Handler {
	// limiter.New() returns a new limiter middleware with the specified configuration.
	return limiter.New(limiter.Config{
		// Max is the maximum number of requests that can be made in the given time frame.
		Max: 60,
		// Expiration is the time frame in which the requests are counted.
		Expiration: 1 * time.Minute,
		// LimiterMiddleware is the storage for the limiter.
		LimiterMiddleware: limiter.SlidingWindow{},
		// LimitReached is a function that is called when the limit is reached.
		LimitReached: func(c *fiber.Ctx) error {
			// response.TooManyRequests() sends a 429 Too Many Requests response.
			return response.TooManyRequests(c, "Too many requests, please try again after one minute.")
		},
		// Next is a function that determines whether to skip this middleware.
		Next: func(c *fiber.Ctx) bool {
			// The middleware is skipped if the request is coming from the server itself.
			return c.IP() == cfg.Server.Host
		},
	})
}

// StrictSecurityLimiter is a middleware that provides strict rate limiting for security-sensitive endpoints.
// It takes the application configuration as input and returns a Fiber handler.
//
// @param cfg *config.Config - The application configuration.
// @return fiber.Handler - The Fiber handler.
func StrictSecurityLimiter(cfg *config.Config) fiber.Handler {
	// limiter.New() returns a new limiter middleware with the specified configuration.
	return limiter.New(limiter.Config{
		// Max is the maximum number of requests that can be made in the given time frame.
		Max: 5,
		// Expiration is the time frame in which the requests are counted.
		Expiration: 10 * time.Minute,
		// LimitReached is a function that is called when the limit is reached.
		LimitReached: func(c *fiber.Ctx) error {
			// response.TooManyRequests() sends a 429 Too Many Requests response.
			return response.TooManyRequests(c, "Too many failed attempts. This action is blocked for 10 minutes.")
		},
		// Next is a function that determines whether to skip this middleware.
		Next: func(c *fiber.Ctx) bool {
			// The middleware is skipped if the request is coming from the server itself.
			return c.IP() == cfg.Server.Host
		},
	})
}