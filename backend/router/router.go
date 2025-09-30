package router

import (
	"database/sql"

	"github.com/rahulcodepython/todo-backend/apps/users"
	"github.com/rahulcodepython/todo-backend/backend/database"
	"github.com/rahulcodepython/todo-backend/backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/rahulcodepython/todo-backend/backend/config"
)

func Router(app *fiber.App, cfg *config.Config, db *sql.DB) {
	app.Use(middleware.Cors(cfg))
	app.Use(middleware.Logger(cfg))

	api := app.Group("/api/v1") // /api/v1

	api.Get("/", func(c *fiber.Ctx) error {
		database.PingDB(db)

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Database connected successfully",
		})
	})

	auth := api.Group("/auth") // /api/v1/auth

	userController := users.NewUserControl(cfg, db)
	auth.Post("/register", userController.RegisterUserController) // /api/v1/auth/register
}
