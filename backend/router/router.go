// router.go
package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/rahulcodepython/todo-backend/apps/users"
	"github.com/rahulcodepython/todo-backend/backend/config"
	"github.com/rahulcodepython/todo-backend/backend/database"
	"github.com/rahulcodepython/todo-backend/backend/middleware"
	"github.com/rahulcodepython/todo-backend/backend/response"
)

func Router(app *fiber.App, cfg *config.Config, db *sql.DB) {
	// Apply global middleware
	app.Use(middleware.Cors(cfg))
	app.Use(middleware.Logger(cfg))

	// Create API group
	api := app.Group("/api/v1")

	// Health check route
	api.Get("/", func(c *fiber.Ctx) error {
		database.PingDB(db)
		return response.OKResponse(c, "Database connected successfully", nil)
	})

	// Create auth group
	auth := api.Group("/auth")

	// Initialize user controller
	userController := users.NewUserControl(cfg, db)

	// Public routes
	auth.Post("/register", userController.RegisterUserController)
	auth.Post("/login", userController.LoginUserController)

	// CRITICAL: Initialize middleware ONCE before using it
	authMiddleware := middleware.Authenticated(db)
	// Protected routes
	auth.Get("/logout", authMiddleware, userController.LogoutUserController)
	auth.Get("/profile", authMiddleware, userController.UserProfileController)
}
