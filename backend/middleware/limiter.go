package middleware

import (
	"time" // Import the "time" package to handle time-related durations for rate limiting.

	"github.com/gofiber/fiber/v2"                            // Import the Fiber web framework, which provides the core functionalities for building web applications in Go.
	"github.com/gofiber/fiber/v2/middleware/limiter"         // Import the Fiber rate limiter middleware, which controls the rate of incoming requests.
	"github.com/rahulcodepython/todo-backend/backend/config" // Import the application's configuration package to access server-related settings.
)

// GeneralAPILimiter creates a rate limiting middleware for general API endpoints.
// This middleware is designed to prevent abuse and ensure fair usage of the API
// by restricting the number of requests a client can make within a specified time window.
// It takes the application's configuration as a parameter to potentially customize behavior.
func GeneralAPILimiter(cfg *config.Config) fiber.Handler {
	// Return a new rate limiter middleware instance configured with specific options.
	return limiter.New(limiter.Config{
		Max:        60,              // Max specifies the maximum number of requests allowed within the Expiration window. Here, 60 requests.
		Expiration: 1 * time.Minute, // Expiration defines the duration of the rate limiting window. Here, 1 minute.
		// LimiterMiddleware specifies the rate limiting algorithm to use.
		// SlidingWindow is a common and effective algorithm that tracks requests over a moving time window.
		LimiterMiddleware: limiter.SlidingWindow{},
		// LimitReached is a custom function that is executed when a client exceeds the defined rate limit.
		// It constructs and sends an appropriate HTTP response to the client.
		LimitReached: func(c *fiber.Ctx) error {
			// Set the HTTP status code to 429 (Too Many Requests) to indicate that the client has sent too many requests in a given amount of time.
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "fail",                                                  // A status field indicating the request failed due to rate limiting.
				"message": "Too many requests, please try again after one minute.", // A user-friendly message explaining the rate limit.
			})
		},
		// Next is a function that determines whether the middleware should be skipped for the current request.
		// In this configuration, the rate limiter is skipped if the client's IP address matches the server's host IP.
		// This is useful for internal requests or health checks, preventing them from being rate-limited.
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == cfg.Server.Host
		},
	})
}

// StrictSecurityLimiter creates a more aggressive rate limiting middleware, typically used for sensitive endpoints
// like login attempts, password resets, or other actions prone to brute-force attacks.
// This limiter has a lower request threshold and a longer blocking duration to enhance security.
// It takes the application's configuration as a parameter to potentially customize behavior.
func StrictSecurityLimiter(cfg *config.Config) fiber.Handler {
	// Return a new rate limiter middleware instance configured with strict security options.
	return limiter.New(limiter.Config{
		Max:        5,                // Max specifies a very low maximum number of requests allowed within the Expiration window. Here, 5 requests.
		Expiration: 10 * time.Minute, // Expiration defines a longer blocking duration once the limit is reached. Here, 10 minutes.
		// LimitReached is a custom function that is executed when a client exceeds this strict rate limit.
		// It provides a specific message indicating a security-related block.
		LimitReached: func(c *fiber.Ctx) error {
			// Set the HTTP status code to 429 (Too Many Requests) to indicate that the client has sent too many requests.
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "fail",                                                             // A status field indicating the request failed due to rate limiting.
				"message": "Too many failed attempts. This action is blocked for 10 minutes.", // A specific message for security-related rate limiting.
			})
		},
		// Next is a function that determines whether the middleware should be skipped for the current request.
		// Similar to GeneralAPILimiter, it skips rate limiting for requests originating from the server's host IP.
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == cfg.Server.Host
		},
	})
}
