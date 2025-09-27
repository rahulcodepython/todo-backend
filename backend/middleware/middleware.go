package middleware

import (
	"github.com/rahulcodepython/todo-backend/backend/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupMiddleware(app *fiber.App, cfg *config.Config) {
	// Add Recover middleware to catch panics
	app.Use(recover.New())

	// Add Logger middleware
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	// Add Rate Limiter middleware
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.Middleware.RateLimitMax,
		Expiration: cfg.Middleware.RateLimitWindow,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Too many requests, please try again later.",
			})
		},
	}))
}
