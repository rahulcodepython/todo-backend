// This file defines the main router for the application.
// It sets up all the API routes and applies the necessary middleware.
package router

// "database/sql" provides a generic SQL interface. It is used here to pass the database connection to the controllers.
import (
	"database/sql"

	// "github.com/gofiber/fiber/v2" is a web framework for Go. It is used here to create the router and define the routes.
	"github.com/gofiber/fiber/v2"
	// "github.com/rahulcodepython/todo-backend/apps/todos" is a local package that contains the todo controllers.
	"github.com/rahulcodepython/todo-backend/apps/todos"
	// "github.com/rahulcodepython/todo-backend/apps/users" is a local package that contains the user controllers.
	"github.com/rahulcodepython/todo-backend/apps/users"
	// "github.com/rahulcodepython/todo-backend/backend/config" is a local package that provides access to the application configuration.
	"github.com/rahulcodepython/todo-backend/backend/config"
	// "github.com/rahulcodepython/todo-backend/backend/database" is a local package that provides database-related functions.
	"github.com/rahulcodepython/todo-backend/backend/database"
	// "github.com/rahulcodepython/todo-backend/backend/middleware" is a local package that provides middleware for the application.
	"github.com/rahulcodepython/todo-backend/backend/middleware"
	// "github.com/rahulcodepython/todo-backend/backend/response" is a local package that provides standardized API responses.
	"github.com/rahulcodepython/todo-backend/backend/response"
)

// Router sets up the application's routes.
// It takes the Fiber app, configuration, and database connection as input.
//
// @param app *fiber.App - The Fiber application.
// @param cfg *config.Config - The application configuration.
// @param db *sql.DB - The database connection.
func Router(app *fiber.App, cfg *config.Config, db *sql.DB) {
	// app.Use() applies middleware to all routes.
	// middleware.Cors() is a middleware that handles Cross-Origin Resource Sharing.
	app.Use(middleware.Cors(cfg))
	// middleware.Logger() is a middleware that logs information about each request.
	app.Use(middleware.Logger(cfg))

	// authMiddleware is a middleware that checks if a user is authenticated.
	authMiddleware := middleware.Authenticated(db)
	// authenticatedUserMiddleware is a middleware that retrieves the authenticated user's information.
	authenticatedUserMiddleware := middleware.AuthenticatedUser(db)

	// api is a new group of routes with the prefix "/api/v1".
	api := app.Group("/api/v1")

	// This defines a GET route for the root of the API group.
	// It serves as a health check endpoint.
	api.Get("/", func(c *fiber.Ctx) error {
		// database.PingDB() checks if the database connection is alive.
		database.PingDB(db)
		// response.OKResponse() sends a 200 OK response with a success message.
		return response.OKResponse(c, "Database connected successfully", nil)
	})

	// auth is a new group of routes with the prefix "/auth".
	auth := api.Group("/auth")

	// userController is a new instance of the user controller.
	userController := users.NewUserControl(cfg, db)

	// This defines a POST route for user registration.
	auth.Post("/register", userController.RegisterUserController)
	// This defines a POST route for user login.
	auth.Post("/login", userController.LoginUserController)

	// This defines a GET route for user logout.
	// It is protected by the authMiddleware.
	auth.Get("/logout", authMiddleware, userController.LogoutUserController)
	// This defines a GET route for retrieving the user's profile.
	// It is protected by both the authMiddleware and the authenticatedUserMiddleware.
	auth.Get("/profile", authMiddleware, authenticatedUserMiddleware, userController.UserProfileController)

	// todo is a new group of routes with the prefix "/todos".
	// It is protected by both the authMiddleware and the authenticatedUserMiddleware.
	todo := api.Group("/todos", authMiddleware, authenticatedUserMiddleware)

	// todoController is a new instance of the todo controller.
	todoController := todos.NewTodoControl(cfg, db)

	// This defines a POST route for creating a new todo.
	todo.Post("/create", todoController.CreateTodoController)
	// This defines a GET route for retrieving all todos.
	todo.Get("/list", todoController.GetTodosController)
	// This defines a PUT route for updating a todo.
	todo.Put("/update/:id", todoController.UpdateTodoController)
	// This defines a PATCH route for completing a todo.
	todo.Patch("/complete/:id", todoController.CompleteTodoController)
	// This defines a DELETE route for deleting a todo.
	todo.Delete("/delete/:id", todoController.DeleteTodoController)
}