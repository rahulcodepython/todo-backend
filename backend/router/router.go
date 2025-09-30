package router

import (
	// Import the "database/sql" package to interact with SQL databases.
	"database/sql"

	// Import the "users" package from the application's "apps" directory, which contains user-related functionalities like controllers.
	"github.com/rahulcodepython/todo-backend/apps/users"
	// Import the "database" package from the application's "backend" directory, likely for database connection utilities.
	"github.com/rahulcodepython/todo-backend/backend/database"
	// Import the "middleware" package from the application's "backend" directory, which provides various HTTP middleware functions.
	"github.com/rahulcodepython/todo-backend/backend/middleware"

	// Import the "fiber" package, a fast, unopinionated, and flexible web framework for Go.
	"github.com/gofiber/fiber/v2"
	// Import the "config" package from the application's "backend" directory, used for loading application configurations.
	"github.com/rahulcodepython/todo-backend/backend/config"
)

// Router sets up all the application's routes and applies global middleware.
// It takes a Fiber app instance, the application configuration, and a database connection as parameters.
func Router(app *fiber.App, cfg *config.Config, db *sql.DB) {
	// Apply the CORS middleware to the Fiber app. This middleware handles Cross-Origin Resource Sharing,
	// allowing or restricting resource requests from different origins based on the configuration.
	app.Use(middleware.Cors(cfg))
	// Apply the Logger middleware to the Fiber app. This middleware logs details about incoming HTTP requests,
	// which is useful for debugging and monitoring.
	app.Use(middleware.Logger(cfg))

	// Create a new Fiber group for API version 1 routes. All routes defined within this group
	// will be prefixed with "/api/v1". This helps in versioning the API.
	api := app.Group("/api/v1") // /api/v1

	// Define a GET route for the root of the API group ("/api/v1/").
	// This route is used to check the database connection status.
	api.Get("/", func(c *fiber.Ctx) error {
		// Call PingDB to check if the database connection is alive. This is a health check for the database.
		database.PingDB(db)

		// Return a JSON response indicating successful database connection.
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Database connected successfully",
		})
	})

	// Create a new Fiber group for authentication-related routes. All routes within this group
	// will be prefixed with "/api/v1/auth".
	auth := api.Group("/auth") // /api/v1/auth

	// Initialize a new UserController with the application configuration and database connection.
	// The UserController handles all business logic related to user management.
	userController := users.NewUserControl(cfg, db)
	// Define a POST route for user registration ("/api/v1/auth/register").
	// This route maps to the RegisterUserController method of the userController.
	auth.Post("/register", userController.RegisterUserController) // /api/v1/auth/register
	// Define a POST route for user login ("/api/v1/auth/login").
	// This route maps to the LoginUserController method of the userController.
	auth.Post("/login", userController.LoginUserController) // /api/v1/auth/login
}
