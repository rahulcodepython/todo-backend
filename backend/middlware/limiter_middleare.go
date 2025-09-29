package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func GeneralAPILimiter() fiber.Handler {
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
	})
}

func StrictSecurityLimiter() fiber.Handler {
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
	})
}
