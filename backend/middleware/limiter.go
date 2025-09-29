package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/rahulcodepython/todo-backend/backend/config"
)

func GeneralAPILimiter(cfg *config.Config) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:               60,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "fail",
				"message": "Too many requests, please try again after one minute.",
			})
		},
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == cfg.Server.Host
		},
	})
}

func StrictSecurityLimiter(cfg *config.Config) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:                    5,
		Expiration:             10 * time.Minute,
		SkipSuccessfulRequests: true,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "fail",
				"message": "Too many failed attempts. This action is blocked for 10 minutes.",
			})
		},
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == cfg.Server.Host
		},
	})
}
