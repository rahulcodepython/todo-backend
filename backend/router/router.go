// File: backend/router/router.go

package router

import (
	"database/sql"

	"github.com/rahulcodepython/todo-backend/apps/todos"
	"github.com/rahulcodepython/todo-backend/apps/users"
	"github.com/rahulcodepython/todo-backend/backend/config"
	"github.com/rahulcodepython/todo-backend/backend/middleware"

	"github.com/gofiber/fiber/v2"
)

// Setup configures the application routes and middleware.
func Setup(app *fiber.App, db *sql.DB, cfg *config.Config) {
	// Setup global middleware (Logger, Recover, Rate Limiter)
	middleware.SetupMiddleware(app, cfg)

	// Create the main API group
	api := app.Group("/api/v1")

	// Health check route
	healthController := NewController(db)
	api.Get("/", healthController.HealthCheckUp)

	// --- Register all installed apps here ---
	// This is where we break the cycle. The router knows about the apps.
	// The config does not. We pass the necessary dependencies (router, db, config)
	// down to each app's route registration function.
	users.RegisterRoutes(api, db, cfg)
	todos.RegisterRoutes(api, db, cfg)
	// To add a new app, you would just add its registration line here.
	// e.g., notes.RegisterRoutes(api, db, cfg)

	// 404 Handler for API routes
	api.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Route not found",
		})
	})
}
