package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rahulcodepython/todo-backend/backend/config"
)

func Logger(cfg *config.Config) fiber.Handler {
	return logger.New(logger.Config{
		Format: "${protocol}://${ip}:${port} ${status} - ${method} ${path} [${time}]\n",
	})
}
