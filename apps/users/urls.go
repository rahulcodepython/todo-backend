package users

import (
	"database/sql"

	"github.com/rahulcodepython/todo-backend/backend/config"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes initializes the routes for the users app.
// It now accepts the db connection pool and config struct directly.
func RegisterRoutes(router fiber.Router, db *sql.DB, cfg *config.Config) {
	// Create a new controller instance with the dependencies.
	controller := NewUsersController(db, cfg)

	// Create a group for auth routes
	authGroup := router.Group("/auth")

	// Public routes
	authGroup.Post("/register", controller.Register)
	authGroup.Post("/login", controller.Login)

	// Protected routes that require the auth middleware
	authGroup.Post("/logout", AuthMiddleware(db), controller.Logout)
	authGroup.Get("/profile", AuthMiddleware(db), controller.GetProfile)
}
