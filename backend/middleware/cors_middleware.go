package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CorsRestriction() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "https://your-production-domain.com",
		AllowHeaders: "Origin, Content-Type, Accept",
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
	})
}
