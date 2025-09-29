package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rahulcodepython/todo-backend/backend/config"
)

func Cors(cfg *config.Config) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: cfg.CORS.CorsOrigins,
		AllowHeaders: "Origin, Content-Type, Accept",
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == cfg.Server.Host
		},
	})
}
