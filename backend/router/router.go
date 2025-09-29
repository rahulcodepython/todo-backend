package router

import (
	"database/sql"

	"github.com/rahulcodepython/todo-backend/backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/rahulcodepython/todo-backend/backend/config"
)

func Router(app *fiber.App, cfg *config.Config, db *sql.DB) {
	app.Use(middleware.Cors(cfg))
	app.Use(middleware.Logger(cfg))
	app.Use(middleware.GeneralAPILimiter(cfg))
	app.Use(middleware.StrictSecurityLimiter(cfg))
	app.Use(middleware.Recover(cfg))

	api := app.Group("/api/v1") // /api/v1

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to API v1",
		})
	}) // /api/v1/

}
