// File: apps/todos/urls.go

package todos

import (
	"database/sql"

	"github.com/rahulcodepython/todo-backend/apps/users"
	"github.com/rahulcodepython/todo-backend/backend/config"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes initializes the routes for the todos app.
func RegisterRoutes(router fiber.Router, db *sql.DB, cfg *config.Config) {
	// Create a new controller instance with the database connection.
	controller := NewTodosController(db)

	// Create a group for todo routes and apply the auth middleware to all of them.
	todoGroup := router.Group("/todos", users.AuthMiddleware(db))

	todoGroup.Post("/", controller.CreateTodo)
	todoGroup.Get("/", controller.GetAllTodos)
	todoGroup.Put("/:id", controller.UpdateTodo)
	todoGroup.Patch("/:id/complete", controller.CompleteTodo)
	todoGroup.Delete("/:id", controller.DeleteTodo)
}
