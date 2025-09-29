package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rahulcodepython/todo-backend/backend/config"
)

func Recover(cfg *config.Config) fiber.Handler {
	return recover.New()
}
